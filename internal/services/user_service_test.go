package services

import (
	"errors"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/utils"
	"testing"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	// Data untuk simulasi
	users          map[uint]*models.User
	usersByEmail   map[string]*models.User
	createError    error
	updateError    error
	deleteError    error
	updatePwdError error
	findAllError   error
	findAllUsers   []models.User
	findAllTotal   int64
}

// NewMockUserRepository membuat mock repository baru
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:        make(map[uint]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
}

// SetUser menambahkan user ke mock untuk testing
func (m *MockUserRepository) SetUser(user *models.User) {
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.createError != nil {
		return m.createError
	}
	if m.usersByEmail[user.Email] != nil {
		return errors.New("email already exists")
	}
	// Simulasikan ID auto-increment
	if user.ID == 0 {
		user.ID = uint(len(m.users) + 1)
	}
	user.DibuatPada = time.Now()
	user.DiperbaruiPada = time.Now()
	m.SetUser(user)
	return nil
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	user, exists := m.usersByEmail[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindAll(page, pageSize int, search, peran string, aktif *bool) ([]models.User, int64, error) {
	if m.findAllError != nil {
		return nil, 0, m.findAllError
	}
	return m.findAllUsers, m.findAllTotal, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if m.updateError != nil {
		return m.updateError
	}
	if _, exists := m.users[user.ID]; !exists {
		return gorm.ErrRecordNotFound
	}
	user.DiperbaruiPada = time.Now()
	m.SetUser(user)
	return nil
}

func (m *MockUserRepository) UpdatePassword(id uint, hashedPassword string) error {
	if m.updatePwdError != nil {
		return m.updatePwdError
	}
	user, exists := m.users[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}
	user.Password = hashedPassword
	user.DiperbaruiPada = time.Now()
	return nil
}

func (m *MockUserRepository) Delete(id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, exists := m.users[id]; !exists {
		return gorm.ErrRecordNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) Count(search, peran string, aktif *bool) (int64, error) {
	return int64(len(m.users)), nil
}

// Helper function untuk membuat user service dengan mock
func newTestUserService(mockRepo repositories.UserRepository) *userService {
	return &userService{
		userRepo: mockRepo,
		logger:   zap.NewNop(), // No-op logger untuk testing
	}
}

// Helper function untuk hash password (sama seperti di service)
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ==================== Test GetUserByID ====================

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: Buat user di mock
	hashedPwd, _ := hashPassword("password123")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
		Nama:     "Test User",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	// Execute
	result, err := service.GetUserByID(1)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected user response, got nil")
	}
	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
	if result.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", result.Email)
	}
	if result.Nama != "Test User" {
		t.Errorf("Expected nama Test User, got %s", result.Nama)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Execute: User tidak ada
	result, err := service.GetUserByID(999)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result, got user")
	}
	if err != utils.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

// ==================== Test CreateUser ====================

func TestCreateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	req := dto.CreateUserRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Nama:     "New User",
		Peran:    "kasir",
	}

	// Execute
	result, err := service.CreateUser(req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected user response, got nil")
	}
	if result.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, result.Email)
	}
	if result.Nama != req.Nama {
		t.Errorf("Expected nama %s, got %s", req.Nama, result.Nama)
	}
	if result.Peran != req.Peran {
		t.Errorf("Expected peran %s, got %s", req.Peran, result.Peran)
	}
	if !result.Aktif {
		t.Error("Expected user to be aktif by default")
	}

	// Verify password di-hash (tidak sama dengan plain password)
	user, _ := mockRepo.FindByEmail(req.Email)
	if user.Password == req.Password {
		t.Error("Password should be hashed, not plain text")
	}
}

func TestCreateUser_EmailExists(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User dengan email sudah ada
	hashedPwd, _ := hashPassword("password123")
	existingUser := &models.User{
		ID:       1,
		Email:    "existing@example.com",
		Password: hashedPwd,
		Nama:     "Existing User",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(existingUser)

	req := dto.CreateUserRequest{
		Email:    "existing@example.com", // Email yang sudah ada
		Password: "password123",
		Nama:     "New User",
		Peran:    "kasir",
	}

	// Execute
	result, err := service.CreateUser(req)

	// Assert
	if err == nil {
		t.Error("Expected error for existing email, got nil")
	}
	if result != nil {
		t.Error("Expected nil result when email exists")
	}
	if err != utils.ErrUserExists {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestCreateUser_DatabaseError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.createError = errors.New("database connection failed")
	service := newTestUserService(mockRepo)

	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Nama:     "Test User",
		Peran:    "kasir",
	}

	// Execute
	result, err := service.CreateUser(req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result on database error")
	}
	appErr := utils.GetAppError(err)
	if appErr.Code != utils.ErrCodeDatabaseError {
		t.Errorf("Expected ErrCodeDatabaseError, got %s", appErr.Code)
	}
}

// ==================== Test ListUsers ====================

func TestListUsers_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: Mock FindAll
	mockRepo.findAllUsers = []models.User{
		{ID: 1, Email: "user1@example.com", Nama: "User 1", Peran: "kasir", Aktif: true},
		{ID: 2, Email: "user2@example.com", Nama: "User 2", Peran: "admin_gudang", Aktif: true},
	}
	mockRepo.findAllTotal = 2

	req := dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	// Execute
	result, err := service.ListUsers(req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected list response, got nil")
	}
	if len(result.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(result.Users))
	}
	if result.Pagination.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Pagination.Total)
	}
	if result.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Pagination.Page)
	}
	if result.Pagination.PageSize != 10 {
		t.Errorf("Expected page size 10, got %d", result.Pagination.PageSize)
	}
	if result.Pagination.TotalPages != 1 {
		t.Errorf("Expected total pages 1, got %d", result.Pagination.TotalPages)
	}
}

func TestListUsers_DefaultPagination(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	mockRepo.findAllUsers = []models.User{}
	mockRepo.findAllTotal = 0

	req := dto.ListUsersRequest{
		Page:     0, // Invalid, should default to 1
		PageSize: 0, // Invalid, should default to 10
	}

	// Execute
	result, err := service.ListUsers(req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Pagination.Page != 1 {
		t.Errorf("Expected page to default to 1, got %d", result.Pagination.Page)
	}
	if result.Pagination.PageSize != 10 {
		t.Errorf("Expected page size to default to 10, got %d", result.Pagination.PageSize)
	}
}

func TestListUsers_MaxPageSize(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	mockRepo.findAllUsers = []models.User{}
	mockRepo.findAllTotal = 0

	req := dto.ListUsersRequest{
		Page:     1,
		PageSize: 200, // Exceeds max (100), should be capped
	}

	// Execute
	result, err := service.ListUsers(req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Pagination.PageSize != 100 {
		t.Errorf("Expected page size to be capped at 100, got %d", result.Pagination.PageSize)
	}
}

func TestListUsers_DatabaseError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.findAllError = errors.New("database error")
	service := newTestUserService(mockRepo)

	req := dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	// Execute
	result, err := service.ListUsers(req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result on database error")
	}
	appErr := utils.GetAppError(err)
	if appErr.Code != utils.ErrCodeDatabaseError {
		t.Errorf("Expected ErrCodeDatabaseError, got %s", appErr.Code)
	}
}

// ==================== Test UpdateUser ====================

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User yang akan di-update
	hashedPwd, _ := hashPassword("password123")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
		Nama:     "Old Name",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	namaBaru := "New Name"
	req := dto.UpdateUserRequest{
		Nama: &namaBaru,
	}

	// Execute
	result, err := service.UpdateUser(1, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected user response, got nil")
	}
	if result.Nama != "New Name" {
		t.Errorf("Expected nama New Name, got %s", result.Nama)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	namaBaru := "New Name"
	req := dto.UpdateUserRequest{
		Nama: &namaBaru,
	}

	// Execute: User tidak ada
	result, err := service.UpdateUser(999, req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result, got user")
	}
	if err != utils.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUser_UpdateAllFields(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup
	hashedPwd, _ := hashPassword("password123")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
		Nama:     "Old Name",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	namaBaru := "New Name"
	peranBaru := "admin_gudang"
	aktifBaru := false
	req := dto.UpdateUserRequest{
		Nama:  &namaBaru,
		Peran: &peranBaru,
		Aktif: &aktifBaru,
	}

	// Execute
	result, err := service.UpdateUser(1, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Nama != "New Name" {
		t.Errorf("Expected nama New Name, got %s", result.Nama)
	}
	if result.Peran != "admin_gudang" {
		t.Errorf("Expected peran admin_gudang, got %s", result.Peran)
	}
	if result.Aktif {
		t.Error("Expected aktif to be false")
	}
}

// ==================== Test DeleteUser ====================

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User yang akan dihapus
	hashedPwd, _ := hashPassword("password123")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
		Nama:     "Test User",
		Peran:    "kasir", // Bukan owner, bisa dihapus
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	// Execute
	err := service.DeleteUser(1)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Verify user sudah dihapus
	_, err = mockRepo.FindByID(1)
	if err == nil {
		t.Error("Expected user to be deleted")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Execute: User tidak ada
	err := service.DeleteUser(999)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != utils.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteUser_CannotDeleteOwner(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User dengan role owner
	hashedPwd, _ := hashPassword("password123")
	owner := &models.User{
		ID:       1,
		Email:    "owner@example.com",
		Password: hashedPwd,
		Nama:     "Owner",
		Peran:    "owner", // Owner tidak bisa dihapus
		Aktif:    true,
	}
	mockRepo.SetUser(owner)

	// Execute
	err := service.DeleteUser(1)

	// Assert
	if err == nil {
		t.Error("Expected error when deleting owner, got nil")
	}
	appErr := utils.GetAppError(err)
	if appErr.Code != utils.ErrCodeValidationError {
		t.Errorf("Expected ErrCodeValidationError, got %s", appErr.Code)
	}
	if appErr.Message != "Cannot delete owner user" {
		t.Errorf("Expected message 'Cannot delete owner user', got %s", appErr.Message)
	}
}

// ==================== Test ChangePassword ====================

func TestChangePassword_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User dengan password lama
	oldPassword := "oldpassword123"
	hashedOldPwd, _ := hashPassword(oldPassword)
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedOldPwd,
		Nama:     "Test User",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	req := dto.ChangePasswordRequest{
		PasswordLama: oldPassword,
		PasswordBaru: "newpassword123",
	}

	// Execute
	err := service.ChangePassword(1, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Verify password sudah di-update
	updatedUser, _ := mockRepo.FindByID(1)
	if updatedUser.Password == hashedOldPwd {
		t.Error("Password should be updated")
	}
	// Verify password baru bisa di-verify
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123"))
	if err != nil {
		t.Error("New password should be correctly hashed")
	}
}

func TestChangePassword_UserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	req := dto.ChangePasswordRequest{
		PasswordLama: "oldpassword123",
		PasswordBaru: "newpassword123",
	}

	// Execute: User tidak ada
	err := service.ChangePassword(999, req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != utils.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestChangePassword_InvalidOldPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	// Setup: User dengan password tertentu
	hashedPwd, _ := hashPassword("correctpassword")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
		Nama:     "Test User",
		Peran:    "kasir",
		Aktif:    true,
	}
	mockRepo.SetUser(user)

	req := dto.ChangePasswordRequest{
		PasswordLama: "wrongpassword", // Password salah
		PasswordBaru: "newpassword123",
	}

	// Execute
	err := service.ChangePassword(1, req)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid old password, got nil")
	}
	appErr := utils.GetAppError(err)
	if appErr.Code != utils.ErrCodeValidationError {
		t.Errorf("Expected ErrCodeValidationError, got %s", appErr.Code)
	}
	if appErr.Message != "Invalid old password" {
		t.Errorf("Expected message 'Invalid old password', got %s", appErr.Message)
	}
}

// ==================== Test toUserResponse ====================

func TestToUserResponse(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := newTestUserService(mockRepo)

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
		Nama:     "Test User",
		Peran:    "kasir",
		Aktif:    true,
	}

	// Execute
	result := service.toUserResponse(user)

	// Assert
	if result == nil {
		t.Fatal("Expected user response, got nil")
	}
	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
	if result.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", result.Email)
	}
	if result.Nama != "Test User" {
		t.Errorf("Expected nama Test User, got %s", result.Nama)
	}
	if result.Peran != "kasir" {
		t.Errorf("Expected peran kasir, got %s", result.Peran)
	}
	if !result.Aktif {
		t.Error("Expected aktif to be true")
	}
	// Verify password tidak ada di response
	// (Password field tidak ada di UserResponse DTO)
}
