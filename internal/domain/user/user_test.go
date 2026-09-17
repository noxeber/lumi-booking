package user

import (
	"errors"
	"testing"
	"uuid"
)

func TestCreateUser(t *testing.T) {
	t.Run("creates user with valid name", func(t *testing.T) {
		name := "Иван Петров"

		user, err := CreateUser(name)
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("CreateUser() returned nil")
		}

		var zeroID uuid.UUID
		if user.id == zeroID {
			t.Error("expected user ID to be generated, got zero UUID")
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("creates user with single character name", func(t *testing.T) {
		name := "А"

		user, err := CreateUser(name)
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("creates user with whitespace-only name", func(t *testing.T) {
		testCases := []string{
			" ",
			"   ",
			"\t",
			"\n",
			" \t\n ",
		}

		for _, name := range testCases {
			t.Run("name="+name, func(t *testing.T) {
				user, err := CreateUser(name)
				if err != nil {
					t.Fatalf("CreateUser(%q) unexpected error: %v", name, err)
				}

				if user.name != name {
					t.Errorf("name = %q, want %q (whitespace should be preserved)", user.name, name)
				}
			})
		}
	})

	t.Run("creates user with long name", func(t *testing.T) {
		name := "Очень длинное имя клиента салона красоты с множеством деталей и особенностей для тестирования"

		user, err := CreateUser(name)
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("creates user with special characters in name", func(t *testing.T) {
		name := "Иван-Петров (младший)"

		user, err := CreateUser(name)
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		user, err := CreateUser("")

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("CreateUser() error = %v, want %v", err, ErrEmptyName)
		}

		if user != nil {
			t.Error("CreateUser() should return nil on error")
		}
	})

	t.Run("generates unique IDs for different users", func(t *testing.T) {
		user1, err := CreateUser("Пользователь 1")
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		user2, err := CreateUser("Пользователь 2")
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if user1.id == user2.id {
			t.Error("expected different IDs for different users")
		}
	})

	t.Run("generates non-zero UUID", func(t *testing.T) {
		user, err := CreateUser("Тестовый пользователь")
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		var zeroID uuid.UUID
		if user.id == zeroID {
			t.Error("expected non-zero UUID, got zero UUID")
		}
	})
}

func TestRestoreUser(t *testing.T) {
	t.Run("restores user with given ID", func(t *testing.T) {
		id := uuid.NewV4()
		name := "Восстановленный пользователь"

		user := RestoreUser(id, name)

		if user == nil {
			t.Fatal("RestoreUser() returned nil")
		}

		if user.id != id {
			t.Errorf("id = %v, want %v", user.id, id)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("does not validate name", func(t *testing.T) {
		id := uuid.NewV4()

		user := RestoreUser(id, "")

		if user == nil {
			t.Fatal("RestoreUser() returned nil")
		}

		if user.name != "" {
			t.Errorf("name = %q, want empty string", user.name)
		}
	})

	t.Run("restores with zero UUID", func(t *testing.T) {
		var zeroID uuid.UUID
		name := "Пользователь с zero ID"

		user := RestoreUser(zeroID, name)

		if user == nil {
			t.Fatal("RestoreUser() returned nil")
		}

		if user.id != zeroID {
			t.Errorf("id = %v, want zero UUID", user.id)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("restores with empty name", func(t *testing.T) {
		id := uuid.NewV4()

		user := RestoreUser(id, "")

		if user.name != "" {
			t.Errorf("name = %q, want empty string", user.name)
		}
	})

	t.Run("restores with whitespace name", func(t *testing.T) {
		id := uuid.NewV4()
		name := "   "

		user := RestoreUser(id, name)

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("preserves all field values exactly", func(t *testing.T) {
		id := uuid.NewV4()
		name := "Точное имя пользователя с особыми символами: !@#$%^&*()"

		user := RestoreUser(id, name)

		if user.id != id {
			t.Errorf("id = %v, want %v", user.id, id)
		}

		if user.name != name {
			t.Errorf("name = %q, want %q", user.name, name)
		}
	})

	t.Run("restores multiple users independently", func(t *testing.T) {
		id1 := uuid.NewV4()
		id2 := uuid.NewV4()
		name1 := "Пользователь 1"
		name2 := "Пользователь 2"

		user1 := RestoreUser(id1, name1)
		user2 := RestoreUser(id2, name2)

		if user1.id == user2.id {
			t.Error("restored users should have different IDs")
		}

		if user1.name == user2.name {
			t.Error("restored users should have different names")
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("ErrEmptyName has correct message", func(t *testing.T) {
		expected := "client name is empty"
		if ErrEmptyName.Error() != expected {
			t.Errorf("ErrEmptyName.Error() = %q, want %q", ErrEmptyName.Error(), expected)
		}
	})
}

func TestUserWorkflow(t *testing.T) {
	t.Run("create and restore workflow", func(t *testing.T) {
		name := "Мария Сидорова"

		// Create user
		created, err := CreateUser(name)
		if err != nil {
			t.Fatalf("CreateUser() unexpected error: %v", err)
		}

		if created.name != name {
			t.Errorf("created.name = %q, want %q", created.name, name)
		}

		var zeroID uuid.UUID
		if created.id == zeroID {
			t.Fatal("created.id should not be zero UUID")
		}

		// Simulate saving to DB and restoring
		restored := RestoreUser(created.id, created.name)

		// Verify all fields match
		if restored.id != created.id {
			t.Errorf("restored.id = %v, want %v", restored.id, created.id)
		}

		if restored.name != created.name {
			t.Errorf("restored.name = %q, want %q", restored.name, created.name)
		}
	})

	t.Run("create multiple users workflow", func(t *testing.T) {
		names := []string{
			"Анна Иванова",
			"Петр Сидоров",
			"Ольга Петрова",
		}

		users := make([]*User, 0, len(names))

		for _, name := range names {
			user, err := CreateUser(name)
			if err != nil {
				t.Fatalf("CreateUser(%q) unexpected error: %v", name, err)
			}
			users = append(users, user)
		}

		// Verify all users have unique IDs
		seenIDs := make(map[uuid.UUID]bool)
		for _, user := range users {
			if seenIDs[user.id] {
				t.Errorf("duplicate ID found: %v", user.id)
			}
			seenIDs[user.id] = true
		}

		// Verify all names are preserved
		for i, user := range users {
			if user.name != names[i] {
				t.Errorf("user[%d].name = %q, want %q", i, user.name, names[i])
			}
		}
	})
}
