package todoitem

import "context"

// deletes a todo the caller owns: a student their own, a counselor one for an assigned student.
// Also deletes access to any media linked to the todo
func (h *Handler) DeleteTodoItem(ctx context.Context, id string) error {
	studentID, err := h.TodoItemRepository.GetTodoItemStudentID(ctx, id)
	if err != nil {
		return err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return err
	}

	err = h.TodoItemRepository.DeleteTodoItem(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
