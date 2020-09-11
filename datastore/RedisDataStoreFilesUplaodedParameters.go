package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Connectors/RedisConnector"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"reflect"
)

type RedisDataStoreFilesUplaodedParameters struct {
}

//GetTypeName returns the type name for a struct
func GetTypeName(myStruct interface{}) string {
	t := reflect.TypeOf(myStruct)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (rds RedisDataStoreFilesUplaodedParameters) Add(item Entities.FilesUploadedParameters) (string, error) {

	//we generate an id for the item
	if item.GetEntityId() == "0" {
		item.Id = shared.GenerateUniqueId(GetTypeName(item) + "-")
	}

	//
	var key = GetTypeName(item)

	//we fetch all the items from the db
	jsonValue, getErr := RedisConnector.GetItemFromRedis(key)

	if getErr != nil {
		return "", getErr
	}

	var allItems []Entities.FilesUploadedParameters

	unmarshallErr := json.Unmarshal([]byte(jsonValue), &allItems)

	if unmarshallErr != nil {
		return "", unmarshallErr
	}

	//we append our item to exisiting items
	allItems = append(allItems, item)

	jsonBytes, marshallErr := json.Marshal(allItems)

	if marshallErr != nil {
		return "", marshallErr
	}

	rawJsonValue := string(jsonBytes)

	//we save items to db
	return item.GetEntityId(), RedisConnector.SaveItemInRedis(key, rawJsonValue)
}

//func ArrayToJson(allItems []Entities.FilesUploadedParameters)

func (rds RedisDataStoreFilesUplaodedParameters) Update(item Entities.FilesUploadedParameters) (string, error) {

	var key = GetTypeName(item)

	//first fetch all items

	//we fetch all the items from the db
	jsonValue, getErr := RedisConnector.GetItemFromRedis(key)

	//no items found
	if getErr != nil {
		return "", getErr
	}

	var allItems []Entities.FilesUploadedParameters

	unmarshallErr := json.Unmarshal([]byte(jsonValue), &allItems)

	if unmarshallErr != nil {
		return "", unmarshallErr
	}

	//items exist in db
	//so we look for that specific item
	isFound, index := checkIfItemWithIdExistsInArray(allItems, item.Id)

	//but no item with the Id specified
	//this is what is debateable...we create the item
	if !isFound {
		return rds.Add(item)
	}

	//items exist in db
	//and item with supplied ID found
	//we replace old item with new one
	allItems[index] = item

	jsonBytes, marshallErr := json.Marshal(allItems)

	if marshallErr != nil {
		return "", marshallErr
	}

	rawJsonValue := string(jsonBytes)

	//save the item
	return item.GetEntityId(), RedisConnector.SaveItemInRedis(key, rawJsonValue)
}

func (rds RedisDataStoreFilesUplaodedParameters) GetById(itemId string) (Entities.FilesUploadedParameters, error) {

	var result = Entities.FilesUploadedParameters{}

	var key = GetTypeName(result)

	//first fetch all items

	//we fetch all the items from the db
	jsonValue, getErr := RedisConnector.GetItemFromRedis(key)

	//no items found
	if getErr != nil {
		return result, getErr
	}

	var allItems []Entities.FilesUploadedParameters

	unmarshallErr := json.Unmarshal([]byte(jsonValue), &allItems)

	if unmarshallErr != nil {
		return result, unmarshallErr
	}

	//items exist in db
	isFound, index := checkIfItemWithIdExistsInArray(allItems, itemId)

	//but no item with the Id specified
	if !isFound {
		return result, errors.New(fmt.Sprintf("Item with Id [%v] Not found", itemId))
	}

	//item found successfully
	return allItems[index], nil
}

func checkIfItemWithIdExistsInArray(items []Entities.FilesUploadedParameters, itemId string) (bool, int) {
	for i, item := range items {
		if item.GetEntityId() == itemId {
			return true, i
		}
	}

	return false, -1
}

func (rds RedisDataStoreFilesUplaodedParameters) GetAll(afterFetchFilter func(item Entities.FilesUploadedParameters) bool, dataStoreQuery string) ([]Entities.FilesUploadedParameters, error) {

	var result = []Entities.FilesUploadedParameters{}

	//get the group key
	var key = GetTypeName(Entities.FilesUploadedParameters{})

	//we fetch all the items from the db
	jsonValue, getErr := RedisConnector.GetItemFromRedis(key)

	//no items found
	if getErr != nil {
		return result, getErr
	}

	var allItems []Entities.FilesUploadedParameters

	unmarshallErr := json.Unmarshal([]byte(jsonValue), &allItems)

	if unmarshallErr != nil {
		return result, unmarshallErr
	}

	var filteredItems []Entities.FilesUploadedParameters

	for _, item := range allItems {
		if afterFetchFilter(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}
