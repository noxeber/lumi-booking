package master

import (
	"errors"
	"testing"
	"uuid"
)

func TestCreateMaster(t *testing.T) {
	t.Run("creates master with valid name", func(t *testing.T) {
		name := "Иван Иванов"

		master, err := CreateMaster(name)
		if err != nil {
			t.Fatalf("CreateMaster() unexpected error: %v", err)
		}

		if master == nil {
			t.Fatal("CreateMaster() returned nil")
		}

		var zeroID uuid.UUID
		if master.id == zeroID {
			t.Error("expected master ID to be generated, got zero UUID")
		}

		if master.name != name {
			t.Errorf("name = %q, want %q", master.name, name)
		}

		if len(master.services) != 0 {
			t.Errorf("services length = %d, want 0", len(master.services))
		}
	})

	t.Run("creates master with spaces in name", func(t *testing.T) {
		name := "  Мария Петрова  "

		master, err := CreateMaster(name)
		if err != nil {
			t.Fatalf("CreateMaster() unexpected error: %v", err)
		}

		if master.name != name {
			t.Errorf("name = %q, want %q (spaces should be preserved)", master.name, name)
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		master, err := CreateMaster("")

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("CreateMaster() error = %v, want %v", err, ErrEmptyName)
		}

		if master != nil {
			t.Error("CreateMaster() should return nil on error")
		}
	})

	t.Run("returns error for whitespace-only name", func(t *testing.T) {
		testCases := []string{
			" ",
			"   ",
			"\t",
			"\n",
			" \t\n ",
		}

		for _, name := range testCases {
			master, err := CreateMaster(name)

			if !errors.Is(err, ErrEmptyName) {
				t.Errorf("CreateMaster(%q) error = %v, want %v", name, err, ErrEmptyName)
			}

			if master != nil {
				t.Errorf("CreateMaster(%q) should return nil on error", name)
			}
		}
	})

	t.Run("generates unique IDs for different masters", func(t *testing.T) {
		master1, err := CreateMaster("Мастер 1")
		if err != nil {
			t.Fatalf("CreateMaster() unexpected error: %v", err)
		}

		master2, err := CreateMaster("Мастер 2")
		if err != nil {
			t.Fatalf("CreateMaster() unexpected error: %v", err)
		}

		if master1.id == master2.id {
			t.Error("expected different IDs for different masters")
		}
	})
}

func TestAddService(t *testing.T) {
	t.Run("adds new service successfully", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID := uuid.New()

		err := master.AddService(serviceID)
		if err != nil {
			t.Fatalf("AddService() unexpected error: %v", err)
		}

		if len(master.services) != 1 {
			t.Fatalf("services length = %d, want 1", len(master.services))
		}

		if master.services[0] != serviceID {
			t.Errorf("services[0] = %v, want %v", master.services[0], serviceID)
		}
	})

	t.Run("adds multiple different services", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}

		for _, serviceID := range serviceIDs {
			err := master.AddService(serviceID)
			if err != nil {
				t.Fatalf("AddService(%v) unexpected error: %v", serviceID, err)
			}
		}

		if len(master.services) != len(serviceIDs) {
			t.Fatalf("services length = %d, want %d", len(master.services), len(serviceIDs))
		}

		for _, serviceID := range serviceIDs {
			found := false
			for _, s := range master.services {
				if s == serviceID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("services should contain %v", serviceID)
			}
		}
	})

	t.Run("returns error when adding duplicate service", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID := uuid.New()

		err := master.AddService(serviceID)
		if err != nil {
			t.Fatalf("first AddService() unexpected error: %v", err)
		}

		err = master.AddService(serviceID)
		if !errors.Is(err, ErrServiceIDIsAlreadyExist) {
			t.Errorf("AddService() error = %v, want %v", err, ErrServiceIDIsAlreadyExist)
		}

		if len(master.services) != 1 {
			t.Errorf("services length = %d, want 1 (should not add duplicate)", len(master.services))
		}
	})

	t.Run("detects duplicate among multiple services", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID1 := uuid.New()
		serviceID2 := uuid.New()
		serviceID3 := uuid.New()

		master.AddService(serviceID1)
		master.AddService(serviceID2)

		err := master.AddService(serviceID2)
		if !errors.Is(err, ErrServiceIDIsAlreadyExist) {
			t.Errorf("AddService() error = %v, want %v", err, ErrServiceIDIsAlreadyExist)
		}

		if len(master.services) != 2 {
			t.Errorf("services length = %d, want 2", len(master.services))
		}

		err = master.AddService(serviceID3)
		if err != nil {
			t.Fatalf("AddService(%v) unexpected error: %v", serviceID3, err)
		}

		if len(master.services) != 3 {
			t.Errorf("services length = %d, want 3", len(master.services))
		}
	})

	t.Run("treats zero UUID as regular value", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		var zeroID uuid.UUID

		err := master.AddService(zeroID)
		if err != nil {
			t.Fatalf("AddService(zero) unexpected error: %v", err)
		}

		err = master.AddService(zeroID)
		if !errors.Is(err, ErrServiceIDIsAlreadyExist) {
			t.Errorf("AddService(zero) second call error = %v, want %v", err, ErrServiceIDIsAlreadyExist)
		}
	})
}

func TestRemoveService(t *testing.T) {
	t.Run("removes existing service", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID := uuid.New()

		master.AddService(serviceID)
		if len(master.services) != 1 {
			t.Fatalf("setup failed: services length = %d, want 1", len(master.services))
		}

		master.RemoveService(serviceID)

		if len(master.services) != 0 {
			t.Errorf("services length = %d, want 0", len(master.services))
		}
	})

	t.Run("removes service from multiple services", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID1 := uuid.New()
		serviceID2 := uuid.New()
		serviceID3 := uuid.New()

		master.AddService(serviceID1)
		master.AddService(serviceID2)
		master.AddService(serviceID3)

		master.RemoveService(serviceID2)

		if len(master.services) != 2 {
			t.Fatalf("services length = %d, want 2", len(master.services))
		}

		for _, s := range master.services {
			if s == serviceID2 {
				t.Error("services should not contain removed service")
			}
		}

		found1, found3 := false, false
		for _, s := range master.services {
			if s == serviceID1 {
				found1 = true
			}
			if s == serviceID3 {
				found3 = true
			}
		}

		if !found1 {
			t.Error("services should still contain serviceID1")
		}
		if !found3 {
			t.Error("services should still contain serviceID3")
		}
	})

	t.Run("does nothing when removing non-existent service", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		existingID := uuid.New()
		nonExistentID := uuid.New()

		master.AddService(existingID)

		master.RemoveService(nonExistentID)

		if len(master.services) != 1 {
			t.Errorf("services length = %d, want 1", len(master.services))
		}

		if master.services[0] != existingID {
			t.Error("existing service should remain unchanged")
		}
	})

	t.Run("does nothing when removing from empty list", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID := uuid.New()

		master.RemoveService(serviceID)

		if len(master.services) != 0 {
			t.Errorf("services length = %d, want 0", len(master.services))
		}
	})

	t.Run("can remove and re-add same service", func(t *testing.T) {
		master, _ := CreateMaster("Тестовый мастер")
		serviceID := uuid.New()

		master.AddService(serviceID)
		master.RemoveService(serviceID)

		if len(master.services) != 0 {
			t.Fatalf("services length = %d, want 0 after removal", len(master.services))
		}

		err := master.AddService(serviceID)
		if err != nil {
			t.Fatalf("AddService() after removal unexpected error: %v", err)
		}

		if len(master.services) != 1 {
			t.Errorf("services length = %d, want 1", len(master.services))
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("ErrEmptyName has correct message", func(t *testing.T) {
		expected := "master name is empty"
		if ErrEmptyName.Error() != expected {
			t.Errorf("ErrEmptyName.Error() = %q, want %q", ErrEmptyName.Error(), expected)
		}
	})

	t.Run("ErrServiceIDIsAlreadyExist has correct message", func(t *testing.T) {
		expected := "master already provides this service"
		if ErrServiceIDIsAlreadyExist.Error() != expected {
			t.Errorf("ErrServiceIDIsAlreadyExist.Error() = %q, want %q", ErrServiceIDIsAlreadyExist.Error(), expected)
		}
	})
}

func TestMasterWorkflow(t *testing.T) {
	t.Run("complete master lifecycle", func(t *testing.T) {
		name := "Анна Сидорова"

		master, err := CreateMaster(name)
		if err != nil {
			t.Fatalf("CreateMaster() unexpected error: %v", err)
		}

		if master.name != name {
			t.Errorf("name = %q, want %q", master.name, name)
		}

		serviceID1 := uuid.New()
		serviceID2 := uuid.New()
		serviceID3 := uuid.New()

		err = master.AddService(serviceID1)
		if err != nil {
			t.Fatalf("AddService(serviceID1) unexpected error: %v", err)
		}

		err = master.AddService(serviceID2)
		if err != nil {
			t.Fatalf("AddService(serviceID2) unexpected error: %v", err)
		}

		err = master.AddService(serviceID1)
		if !errors.Is(err, ErrServiceIDIsAlreadyExist) {
			t.Fatalf("AddService(duplicate) error = %v, want %v", err, ErrServiceIDIsAlreadyExist)
		}

		if len(master.services) != 2 {
			t.Fatalf("services length = %d, want 2", len(master.services))
		}

		master.RemoveService(serviceID1)
		if len(master.services) != 1 {
			t.Fatalf("services length = %d, want 1 after removal", len(master.services))
		}

		err = master.AddService(serviceID3)
		if err != nil {
			t.Fatalf("AddService(serviceID3) unexpected error: %v", err)
		}

		if len(master.services) != 2 {
			t.Errorf("services length = %d, want 2", len(master.services))
		}

		found2, found3 := false, false
		for _, s := range master.services {
			if s == serviceID2 {
				found2 = true
			}
			if s == serviceID3 {
				found3 = true
			}
		}

		if !found2 || !found3 {
			t.Error("services should contain serviceID2 and serviceID3")
		}
	})
}
