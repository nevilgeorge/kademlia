package kademlia

import (
	"fmt"
	"sync"
)

// KBucket struct
type KBucket struct {
	NodeID           ID
	ContactList      []Contact
	ContactMutexLock sync.Mutex
}

// Initialize KBuckets, called 160 times when Kademlia is instantiated in kademlia.go
func (kb *KBucket) Initialize() {
	kb.NodeID = NewRandomID()
	// create slice for ContactList
	kb.ContactList = make([]Contact, 0, k)
}

// Remove the contact corresponding to a given ID from the KBucket
func (kb *KBucket) RemoveContact(targetID ID) bool {
	fmt.Println("RemoveContact")
	for i, _ := range kb.ContactList {
		if kb.ContactList[i].NodeID == targetID {
			temp := kb.ContactList
			a := append(temp[:i], temp[(i+1):]...)
			kb.ContactList = a
			return true
		}
	}
	return false
}

// Adds a given contact to the end of the kbucket
// Assumes that it is already thread safe
func (kb *KBucket) AddContact(contactList *[]Contact, newContact Contact) {
	fmt.Println("AddContact")
	toAdd := new(Contact)
	toAdd.NodeID = newContact.NodeID
	toAdd.Host = newContact.Host
	toAdd.Port = newContact.Port
	*contactList = append(*contactList, *toAdd)
	fmt.Println("Now I have ...")
	fmt.Println(len(kb.ContactList))
}

// returns a boolean for whether a given Contact exists in the KBucket and index if it was found
func (kb *KBucket) ContainsContact(cont Contact) (exists bool, index int) {
	// iterate through ContactList and compare Contact NodeIDs
	for i := 0; i < len(kb.ContactList); i++ {
		current := kb.ContactList[i]
		if current.NodeID.Equals(cont.NodeID) {
			exists = true
			index = i
			return exists, i
		}
	}
	if !exists {
		exists = false
		index = -1
	}
	return exists, index
}

// Update the KBucket to sort the nodes with most recently used in at the head of the KBucket
func (kb *KBucket) Update(updated Contact) {
	fmt.Println("Update")
	fmt.Println(len(kb.ContactList))
	fmt.Print("K is : ")
	fmt.Print(k)

	// check whether the updated contact exists in the KBucket
	exists, _ := kb.ContainsContact(updated)
	if exists {
		fmt.Println("It exists!")
		// move Contact to the end of the KBucket
		kb.MoveToTail(updated)
	} else if len(kb.ContactList) < k {
		fmt.Println("New contact")
		// create a new contact for the node and add it to the tail of the KBucket
		// not sure if a new Contact needs to be created, but that's what the doc says
		//temp := Contact(CopyID(updated.NodeID), updated.Host, updated.Port)
		temp := new(Contact)
		temp.NodeID = CopyID(updated.NodeID)
		temp.Host = updated.Host
		temp.Port = updated.Port
		fmt.Println("NodeID:")
		fmt.Println(temp.NodeID)
		fmt.Println("Host:")
		fmt.Println(temp.Host)
		fmt.Println("Port:")
		fmt.Println(temp.Port)

		kb.ContactMutexLock.Lock()
		kb.AddContact(&kb.ContactList, *temp) // jwhang: kinda fishy.. not sure if this is ok
		kb.ContactMutexLock.Unlock()
	} else {
		// ping first node in slice
		// if it doesn't respond, removeContact(oldContact) and addContact(updated)
		// else moveToTail(oldContact) and ignore updated
		//firstContact := kb.ContactList[0]
		/*
			ret := kb.Kad.DoPing(firstContact.Host, firstContact.Port)
			if ret == "Error" { // jwhang TODO: Fix this to nil
				kb.ContactMutexLock.Lock()
				kb.RemoveContact(firstContact.NodeID)
				kb.ContactMutexlock.Unlock()

				kb.ContactMutexLock.Lock()
				kb.AddContact(&kb.ContactList, updated)
				kb.ContactMutexLock.Unlock()
			} else {
				kb.MoveToTail(firstContact)
			}
		*/
	}
}

// moves a contact from its position in the KBucket to the end of the same KBucket
func (kb *KBucket) MoveToTail(updated Contact) {
	exists, _ := kb.ContainsContact(updated)
	if exists {
		// finds and removes contact if already exists
		kb.ContactMutexLock.Lock()
		kb.RemoveContact(updated.NodeID)
		kb.ContactMutexLock.Unlock()
	}
	// adds to the end of the KBucket
	kb.ContactMutexLock.Lock()
	kb.AddContact(&kb.ContactList, updated)
	kb.ContactMutexLock.Unlock()
}
