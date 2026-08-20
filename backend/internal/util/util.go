package util

import (
	"fmt"
	"math"
	"time"
)

// GenerateOrderNo generates a unique order number: ORD + timestamp + random 4 digits.
func GenerateOrderNo() string {
	return fmt.Sprintf("ORD%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// GenerateInvoiceNo generates a unique invoice number: INV + timestamp + random 4 digits.
func GenerateInvoiceNo() string {
	return fmt.Sprintf("INV%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// GenerateTicketNo generates a unique ticket number: TIC + timestamp + random 4 digits.
func GenerateTicketNo() string {
	return fmt.Sprintf("TIC%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// GenerateTransactionNo generates a unique transaction number: TXN + timestamp + random 4 digits.
func GenerateTransactionNo() string {
	return fmt.Sprintf("TXN%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// Paginate normalizes page/pageSize and returns (offset, limit).
func Paginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}

// ParseTime attempts to parse a time string in RFC3339 or common date formats.
func ParseTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %s", s)
}
