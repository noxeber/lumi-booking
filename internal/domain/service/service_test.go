package service

import (
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
)

func TestCreateService(t *testing.T) {
	t.Run("creates service with valid parameters", func(t *testing.T) {
		name := "Стрижка"
		price := decimal.NewFromFloat(1500.50)
		duration := 45 * time.Minute
		description := "Мужская и женская стрижка"

		svc, err := CreateService(name, price, duration, description)
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc == nil {
			t.Fatal("CreateService() returned nil")
		}

		var zeroID uuid.UUID
		if svc.id == zeroID {
			t.Error("expected service ID to be generated, got zero UUID")
		}

		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}

		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}

		if svc.duration != duration {
			t.Errorf("duration = %v, want %v", svc.duration, duration)
		}

		if svc.description != description {
			t.Errorf("description = %q, want %q", svc.description, description)
		}
	})

	t.Run("creates service with name exactly 4 characters", func(t *testing.T) {
		name := "Маст" // 4 байта для кириллицы это 2 символа, но len считает байты

		svc, err := CreateService(name, decimal.NewFromInt(100), time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}
	})

	t.Run("creates service with name exactly 4 bytes", func(t *testing.T) {
		name := "Hair"

		svc, err := CreateService(name, decimal.NewFromInt(100), time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}
	})

	t.Run("returns error for name shorter than 4 bytes", func(t *testing.T) {
		testCases := []struct {
			name string
			desc string
		}{
			{"", "empty string"},
			{"a", "1 byte"},
			{"ab", "2 bytes"},
			{"abc", "3 bytes"},
			{"М", "2 bytes (cyrillic)"},
			{"Ма", "4 bytes (cyrillic)"},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				svc, err := CreateService(tc.name, decimal.NewFromInt(100), time.Hour, "")

				if !errors.Is(err, ErrInvalidName) {
					t.Errorf("CreateService(%q) error = %v, want %v", tc.name, err, ErrInvalidName)
				}

				if svc != nil {
					t.Errorf("CreateService(%q) should return nil on error", tc.name)
				}
			})
		}
	})

	t.Run("returns error for negative price", func(t *testing.T) {
		testCases := []decimal.Decimal{
			decimal.NewFromFloat(-0.01),
			decimal.NewFromFloat(-100),
			decimal.NewFromFloat(-999.99),
		}

		for _, price := range testCases {
			svc, err := CreateService("Стрижка", price, time.Hour, "")

			if !errors.Is(err, ErrNegativePrice) {
				t.Errorf("CreateService(price=%v) error = %v, want %v", price, err, ErrNegativePrice)
			}

			if svc != nil {
				t.Errorf("CreateService(price=%v) should return nil on error", price)
			}
		}
	})

	t.Run("accepts zero price", func(t *testing.T) {
		price := decimal.Zero

		svc, err := CreateService("Консультация", price, 15*time.Minute, "Бесплатная консультация")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}
	})

	t.Run("accepts positive price", func(t *testing.T) {
		price := decimal.NewFromFloat(2500.99)

		svc, err := CreateService("Окрашивание", price, 2*time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}
	})

	t.Run("accepts zero duration", func(t *testing.T) {
		duration := time.Duration(0)

		svc, err := CreateService("Тестовая услуга", decimal.NewFromInt(100), duration, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.duration != duration {
			t.Errorf("duration = %v, want %v", svc.duration, duration)
		}
	})

	t.Run("accepts negative duration (no validation)", func(t *testing.T) {
		duration := -1 * time.Hour

		svc, err := CreateService("Тестовая услуга", decimal.NewFromInt(100), duration, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.duration != duration {
			t.Errorf("duration = %v, want %v", svc.duration, duration)
		}
	})

	t.Run("creates service with empty description", func(t *testing.T) {
		svc, err := CreateService("Стрижка", decimal.NewFromInt(1500), time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.description != "" {
			t.Errorf("description = %q, want empty string", svc.description)
		}
	})

	t.Run("creates service with long name", func(t *testing.T) {
		name := "Очень длинное название услуги для салона красоты с множеством деталей и особенностей"

		svc, err := CreateService(name, decimal.NewFromInt(5000), 3*time.Hour, "Подробное описание")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}
	})

	t.Run("creates service with high precision price", func(t *testing.T) {
		price := decimal.RequireFromString("1234.567890123456789")

		svc, err := CreateService("Премиум услуга", price, time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}
	})

	t.Run("generates unique IDs for different services", func(t *testing.T) {
		svc1, err := CreateService("Услуга 1", decimal.NewFromInt(100), time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		svc2, err := CreateService("Услуга 2", decimal.NewFromInt(200), time.Hour, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		if svc1.id == svc2.id {
			t.Error("expected different IDs for different services")
		}
	})

	t.Run("validation errors take precedence", func(t *testing.T) {
		// Both invalid name and negative price
		svc, err := CreateService("ab", decimal.NewFromFloat(-100), time.Hour, "")

		// Should get name error first
		if !errors.Is(err, ErrInvalidName) {
			t.Errorf("CreateService() error = %v, want %v (name validation should be first)", err, ErrInvalidName)
		}

		if svc != nil {
			t.Error("CreateService() should return nil on error")
		}
	})
}

func TestRestoreService(t *testing.T) {
	t.Run("restores service with given ID", func(t *testing.T) {
		id := uuid.NewV4()
		name := "Восстановленная услуга"
		price := decimal.NewFromFloat(2000)
		duration := 90 * time.Minute
		description := "Описание восстановленной услуги"

		svc := RestoreService(id, name, price, duration, description)

		if svc == nil {
			t.Fatal("RestoreService() returned nil")
		}

		if svc.id != id {
			t.Errorf("id = %v, want %v", svc.id, id)
		}

		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}

		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}

		if svc.duration != duration {
			t.Errorf("duration = %v, want %v", svc.duration, duration)
		}

		if svc.description != description {
			t.Errorf("description = %q, want %q", svc.description, description)
		}
	})

	t.Run("does not validate name", func(t *testing.T) {
		id := uuid.NewV4()

		svc := RestoreService(id, "", decimal.NewFromInt(100), time.Hour, "")

		if svc == nil {
			t.Fatal("RestoreService() returned nil")
		}

		if svc.name != "" {
			t.Errorf("name = %q, want empty string", svc.name)
		}

		svc2 := RestoreService(id, "ab", decimal.NewFromInt(100), time.Hour, "")

		if svc2.name != "ab" {
			t.Errorf("name = %q, want %q", svc2.name, "ab")
		}
	})

	t.Run("does not validate price", func(t *testing.T) {
		id := uuid.NewV4()
		negativePrice := decimal.NewFromFloat(-500)

		svc := RestoreService(id, "Услуга", negativePrice, time.Hour, "")

		if svc == nil {
			t.Fatal("RestoreService() returned nil")
		}

		if !svc.price.Equal(negativePrice) {
			t.Errorf("price = %v, want %v", svc.price, negativePrice)
		}
	})

	t.Run("restores with zero UUID", func(t *testing.T) {
		var zeroID uuid.UUID

		svc := RestoreService(zeroID, "Услуга", decimal.NewFromInt(100), time.Hour, "")

		if svc.id != zeroID {
			t.Errorf("id = %v, want zero UUID", svc.id)
		}
	})

	t.Run("restores with empty description", func(t *testing.T) {
		id := uuid.NewV4()

		svc := RestoreService(id, "Услуга", decimal.NewFromInt(100), time.Hour, "")

		if svc.description != "" {
			t.Errorf("description = %q, want empty string", svc.description)
		}
	})

	t.Run("preserves all field values exactly", func(t *testing.T) {
		id := uuid.NewV4()
		name := "Точное имя услуги"
		price := decimal.RequireFromString("999.999999")
		duration := 2*time.Hour + 30*time.Minute + 45*time.Second
		description := "Точное описание с особыми символами: !@#$%^&*()"

		svc := RestoreService(id, name, price, duration, description)

		if svc.id != id {
			t.Errorf("id = %v, want %v", svc.id, id)
		}
		if svc.name != name {
			t.Errorf("name = %q, want %q", svc.name, name)
		}
		if !svc.price.Equal(price) {
			t.Errorf("price = %v, want %v", svc.price, price)
		}
		if svc.duration != duration {
			t.Errorf("duration = %v, want %v", svc.duration, duration)
		}
		if svc.description != description {
			t.Errorf("description = %q, want %q", svc.description, description)
		}
	})
}

func TestDuration(t *testing.T) {
	t.Run("returns correct duration", func(t *testing.T) {
		duration := 90 * time.Minute

		svc, err := CreateService("Тестовая услуга", decimal.NewFromInt(1000), duration, "")
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		got := svc.Duration()

		if got != duration {
			t.Errorf("Duration() = %v, want %v", got, duration)
		}
	})

	t.Run("returns zero duration for zero value", func(t *testing.T) {
		var svc Service

		got := svc.Duration()

		if got != 0 {
			t.Errorf("Duration() = %v, want 0", got)
		}
	})

	t.Run("returns duration from restored service", func(t *testing.T) {
		duration := 2 * time.Hour

		svc := RestoreService(uuid.NewV4(), "Услуга", decimal.NewFromInt(100), duration, "")

		got := svc.Duration()

		if got != duration {
			t.Errorf("Duration() = %v, want %v", got, duration)
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("ErrInvalidName has correct message", func(t *testing.T) {
		expected := "service name is empty or too short"
		if ErrInvalidName.Error() != expected {
			t.Errorf("ErrInvalidName.Error() = %q, want %q", ErrInvalidName.Error(), expected)
		}
	})

	t.Run("ErrNegativePrice has correct message", func(t *testing.T) {
		expected := "service price must be positive"
		if ErrNegativePrice.Error() != expected {
			t.Errorf("ErrNegativePrice.Error() = %q, want %q", ErrNegativePrice.Error(), expected)
		}
	})
}

func TestServiceWorkflow(t *testing.T) {
	t.Run("create and restore workflow", func(t *testing.T) {
		name := "Маникюр с покрытием"
		price := decimal.NewFromFloat(2500)
		duration := 2 * time.Hour
		description := "Классический маникюр с гель-лаком"

		// Create service
		created, err := CreateService(name, price, duration, description)
		if err != nil {
			t.Fatalf("CreateService() unexpected error: %v", err)
		}

		// Simulate saving to DB and restoring
		restored := RestoreService(
			created.id,
			created.name,
			created.price,
			created.duration,
			created.description,
		)

		// Verify all fields match
		if restored.id != created.id {
			t.Errorf("id = %v, want %v", restored.id, created.id)
		}
		if restored.name != created.name {
			t.Errorf("name = %q, want %q", restored.name, created.name)
		}
		if !restored.price.Equal(created.price) {
			t.Errorf("price = %v, want %v", restored.price, created.price)
		}
		if restored.Duration() != created.Duration() {
			t.Errorf("duration = %v, want %v", restored.Duration(), created.Duration())
		}
		if restored.description != created.description {
			t.Errorf("description = %q, want %q", restored.description, created.description)
		}
	})
}
