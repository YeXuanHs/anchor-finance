package backup

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"anchorfinance/pkg/logger"
)

// Service 数据库备份服务
type Service struct {
	backupDir string
	log       *logger.Logger
}

// NewService 创建备份服务
func NewService(backupDir string, log *logger.Logger) *Service {
	// 确保备份目录存在
	os.MkdirAll(backupDir, 0755)
	return &Service{
		backupDir: backupDir,
		log:       log,
	}
}

// BackupConfig 备份配置
type BackupConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// BackupResult 备份结果
type BackupResult struct {
	Success    bool   `json:"success"`
	Filename   string `json:"filename"`
	Filepath   string `json:"filepath"`
	Size       int64  `json:"size"`
	Compressed bool   `json:"compressed"`
	Error      string `json:"error,omitempty"`
}

// Backup 执行数据库备份
func (s *Service) Backup(config BackupConfig) (*BackupResult, error) {
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.sql", config.Database, timestamp)
	filepath := filepath.Join(s.backupDir, filename)

	s.log.Infof("开始备份数据库: %s", config.Database)

	// 执行 mysqldump
	cmd := exec.Command("mysqldump",
		fmt.Sprintf("-h%s", config.Host),
		fmt.Sprintf("-P%d", config.Port),
		fmt.Sprintf("-u%s", config.User),
		fmt.Sprintf("-p%s", config.Password),
		"--single-transaction",
		"--routines",
		"--triggers",
		"--events",
		config.Database,
	)

	outputFile, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("创建备份文件失败: %v", err)
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.Remove(filepath)
		return nil, fmt.Errorf("mysqldump 执行失败: %v", err)
	}

	// 获取文件大小
 fileInfo, _ := os.Stat(filepath)
 size := fileInfo.Size()

	// 压缩备份文件
 compressedPath := filepath + ".gz"
 if err := compressFile(filepath, compressedPath); err != nil {
		s.log.Warnf("压缩备份文件失败: %v", err)
	} else {
		// 删除原始文件
		os.Remove(filepath)
		filepath = compressedPath
		filename = filename + ".gz"
		if info, err := os.Stat(filepath); err == nil {
			size = info.Size()
		}
	}

	s.log.Infof("数据库备份完成: %s (%.2f MB)", filename, float64(size)/1024/1024)

	return &BackupResult{
		Success:    true,
		Filename:   filename,
		Filepath:   filepath,
		Size:       size,
		Compressed: true,
	}, nil
}

// ListBackups 列出所有备份文件
func (s *Service) ListBackups() ([]BackupResult, error) {
	var backups []BackupResult

	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupResult{
			Success:    true,
			Filename:   entry.Name(),
			Filepath:   filepath.Join(s.backupDir, entry.Name()),
			Size:       info.Size(),
			Compressed: filepath.Ext(entry.Name()) == ".gz",
		})
	}

	return backups, nil
}

// DeleteBackup 删除备份文件
func (s *Service) DeleteBackup(filename string) error {
	filepath := filepath.Join(s.backupDir, filename)
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("删除备份文件失败: %v", err)
	}
	s.log.Infof("已删除备份文件: %s", filename)
	return nil
}

// CleanOldBackups 清理旧备份
func (s *Service) CleanOldBackups(keepDays int) (int, error) {
	cutoff := time.Now().AddDate(0, 0, -keepDays)
	deleted := 0

	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(s.backupDir, entry.Name()))
			deleted++
		}
	}

	s.log.Infof("清理了 %d 个过期备份文件", deleted)
	return deleted, nil
}

// Restore 恢复数据库备份
func (s *Service) Restore(config BackupConfig, backupFile string) error {
	filepath := filepath.Join(s.backupDir, backupFile)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupFile)
	}

	// 如果是压缩文件，先解压
	actualFile := filepath
	if filepath[len(filepath)-3:] == ".gz" {
		decompressedFile := filepath[:len(filepath)-3]
		if err := decompressFile(filepath, decompressedFile); err != nil {
			return fmt.Errorf("解压备份文件失败: %v", err)
		}
		defer os.Remove(decompressedFile)
		actualFile = decompressedFile
	}

	s.log.Infof("开始恢复数据库: %s", config.Database)

	// 执行 mysql 恢复
	cmd := exec.Command("mysql",
		fmt.Sprintf("-h%s", config.Host),
		fmt.Sprintf("-P%d", config.Port),
		fmt.Sprintf("-u%s", config.User),
		fmt.Sprintf("-p%s", config.Password),
		config.Database,
	)

	inputFile, err := os.Open(actualFile)
	if err != nil {
		return fmt.Errorf("打开备份文件失败: %v", err)
	}
	defer inputFile.Close()

	cmd.Stdin = inputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql 恢复失败: %v", err)
	}

	s.log.Infof("数据库恢复完成: %s", config.Database)
	return nil
}

// compressFile 压缩文件（gzip）
func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	writer := gzip.NewWriter(dstFile)
	defer writer.Close()

	_, err = io.Copy(writer, srcFile)
	return err
}

// decompressFile 解压文件（gzip）
func decompressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	reader, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, reader)
	return err
}
