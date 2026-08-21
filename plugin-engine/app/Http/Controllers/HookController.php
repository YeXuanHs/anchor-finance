<?php

namespace App\Http\Controllers;

use App\ZjmfCompat\HookCompat;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

/**
 * Hook控制器
 * 处理Go后端发来的Hook触发请求
 */
class HookController extends Controller
{
    /**
     * 触发Hook
     * @param Request $request
     * @return JsonResponse
     */
    public function trigger(Request $request): JsonResponse
    {
        $hook = $request->input('hook');
        $params = $request->input('params', []);

        if (empty($hook)) {
            return response()->json([
                'status' => 'error',
                'message' => 'Hook name is required',
            ], 400);
        }

        // 触发Hook
        $results = HookCompat::trigger($hook, $params);

        return response()->json([
            'status' => 'success',
            'hook' => $hook,
            'results' => $results,
        ]);
    }

    /**
     * 获取Hook列表
     * @param Request $request
     * @return JsonResponse
     */
    public function list(Request $request): JsonResponse
    {
        $hooks = HookCompat::getAllHooks();

        $hookList = [];
        foreach ($hooks as $tag => $handlers) {
            $hookList[] = [
                'tag' => $tag,
                'handler_count' => count($handlers),
            ];
        }

        return response()->json([
            'status' => 'success',
            'data' => $hookList,
        ]);
    }

    /**
     * 注册Hook
     * @param Request $request
     * @return JsonResponse
     */
    public function register(Request $request): JsonResponse
    {
        $tag = $request->input('tag');
        $pluginKey = $request->input('plugin_key');
        $handlerClass = $request->input('handler_class');
        $priority = $request->input('priority', 10);

        if (empty($tag) || empty($handlerClass)) {
            return response()->json([
                'status' => 'error',
                'message' => 'Tag and handler_class are required',
            ], 400);
        }

        if (!class_exists($handlerClass)) {
            return response()->json([
                'status' => 'error',
                'message' => 'Handler class not found: ' . $handlerClass,
            ], 400);
        }

        $handler = [new $handlerClass(), 'handle'];
        HookCompat::add($tag, $handler, $priority);

        return response()->json([
            'status' => 'success',
            'message' => 'Hook registered successfully',
        ]);
    }
}
