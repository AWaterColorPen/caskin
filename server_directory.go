package caskin

// CreateDirectory
// 1. if there exist the object but soft deleted then recover it
// 2. if there does not exist the object then create a new one
func (s *server) CreateDirectory(user User, domain Domain, object Object) error {
	if err := s.RecoverObject(user, domain, object); err == nil {
		return nil
	}
	return s.CreateObject(user, domain, object)
}

// UpdateDirectory
// 1. if there exist the object update properties and g2 in the domain
func (s *server) UpdateDirectory(user User, domain Domain, object Object) error {
	return s.UpdateObject(user, domain, object)
}

// DeleteDirectory
// request.type required as object type
// request.to required as object id
// request.ActionDirectory required for to delete directory item
// if there exist the object
// 1. delete all item of object
// 2. soft delete one object in metadata database
// 3. delete all son of the object in the domain
func (s *server) DeleteDirectory(user User, domain Domain, request *DirectoryRequest) error {
	if request.Type == "" || request.To == 0 || request.ActionDirectory == nil {
		return ErrInValidRequest
	}
	list, err := s.GetObject(user, domain, Write, request.Type)
	if err != nil {
		return err
	}
	var directory []*Directory
	for _, v := range list {
		directory = append(directory, &Directory{Object: v})
	}
	id := ID(NewObjectDirectory(directory).Search(request.To, DirectorySearchAll))
	id = append(id, request.To)
	if err = request.ActionDirectory(id); err != nil {
		return err
	}
	object := DefaultFactory().NewObject()
	object.SetID(request.To)
	return s.DeleteObject(user, domain, object)
}

// GetDirectory
// request.type required as object type
// request.to required as the target object id
// request.CountDirectory required to count the item number of directory
// request.search_type optional as directory search option
//   "top" will only search top level directory by default
//   "all" will search all level directory
// 1. current user has read of object to get directory
func (s *server) GetDirectory(user User, domain Domain, request *DirectoryRequest) ([]*Directory, error) {
	if request.Type == "" || request.To == 0 || request.CountDirectory == nil {
		return nil, ErrInValidRequest
	}
	object, err := s.GetObject(user, domain, Read, request.Type)
	if err != nil {
		return nil, err
	}
	count, err := request.CountDirectory(ID(object))
	if err != nil {
		return nil, err
	}

	var directory []*Directory
	for _, v := range object {
		u := &Directory{Object: v}
		u.TopItemCount = count[u.GetID()]
		directory = append(directory, u)
	}
	return NewObjectDirectory(directory).Search(request.To, request.SearchType), nil
}

// MoveDirectory
// request.to required as the target object id
// request.id required as the source object id list
// request.policy optional as the error handle policy.
//   it will stop by default when error happen
//   "continue" will ignore error
func (s *server) MoveDirectory(user User, domain Domain, request *DirectoryRequest) (*DirectoryResponse, error) {
	if request.Type == "" || request.To == 0 {
		return nil, ErrInValidRequest
	}
	response := &DirectoryResponse{ToDoDirectoryCount: uint64(len(request.ID))}
	object, err := s.DB.GetObjectByID(request.ID)
	if err != nil {
		return nil, err
	}
	for _, v := range object {
		v.SetParentID(request.To)
		if err = s.UpdateObject(user, domain, v); err != nil {
			if request.Policy == "continue" {
				continue
			}
			return response, err
		}
		response.DoneDirectoryCount++
		response.ToDoDirectoryCount--
	}
	return response, nil
}

// MoveItem
// request.to required as the target object id
// request.id required as the source data id list
// request.policy optional as the error handle policy.
//   it will stop by default when error happen
//   "continue" will ignore error
func (s *server) MoveItem(user User, domain Domain, data ObjectData, request *DirectoryRequest) (*DirectoryResponse, error) {
	if request.Type == "" || request.To == 0 {
		return nil, ErrInValidRequest
	}
	response := &DirectoryResponse{ToDoItemCount: uint64(len(request.ID))}
	for _, id := range request.ID {
		item := newByE(data)
		item.SetID(id)
		item.SetObjectID(request.To)
		if err := s.UpdateObjectData(user, domain, item, request.Type); err != nil {
			if request.Policy == "continue" {
				continue
			}
			return response, err
		}
		response.DoneItemCount++
		response.ToDoItemCount--
	}
	return response, nil
}

// CopyItem
// request.to required as the target object id
// request.id required as the source data id list
// request.policy optional as the error handle policy.
//   it will stop by default when error happen
//   "continue" will ignore error
func (s *server) CopyItem(user User, domain Domain, data ObjectData, request *DirectoryRequest) (*DirectoryResponse, error) {
	if request.Type == "" || request.To == 0 {
		return nil, ErrInValidRequest
	}
	response := &DirectoryResponse{ToDoItemCount: uint64(len(request.ID))}
	for _, id := range request.ID {
		item := newByE(data)
		item.SetID(id)
		item.SetDomainID(domain.GetID())
		if err := s.CheckGetObjectData(user, domain, item); err != nil {
			if request.Policy == "continue" {
				continue
			}
			return response, err
		}
		item.SetID(0)
		item.SetObjectID(request.To)
		if err := s.CreateObjectData(user, domain, item, request.Type); err != nil {
			if request.Policy == "continue" {
				continue
			}
			return response, err
		}
		response.DoneItemCount++
		response.ToDoItemCount--
	}
	return response, nil
}
