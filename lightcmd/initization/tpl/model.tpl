type {{ .Model | ucFirst }} {
  id: ID! 
  createdAt: Time!
  updatedAt: Time!
  deletedAt: DeletedAt
  
}
