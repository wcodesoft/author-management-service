package database

import authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"

// AuthorFromGrpc Transforms an Author proto into an Author object.
func AuthorFromGrpc(author *authorManagementProto.Author) Author {
	parsedUuid := uuidParseOrCreate(*author.Uuid)
	return Author{
		ID:     &parsedUuid,
		Name:   author.Name,
		PicURL: author.PicUrl,
	}
}

// AuthorToGrpc Transforms an Author object into a proto Author.
func AuthorToGrpc(author Author) *authorManagementProto.Author {
	uuidString := author.ID.String()
	return &authorManagementProto.Author{
		Uuid:   &uuidString,
		Name:   author.Name,
		PicUrl: author.PicURL,
	}
}

// AuthorListToGrpcList Transforms a list of Author into a AuthorList.
func AuthorListToGrpcList(author []Author) authorManagementProto.AuthorList {
	var parsedAuthors []*authorManagementProto.Author
	for _, author := range author {
		parsedAuthors = append(parsedAuthors, AuthorToGrpc(author))
	}
	return authorManagementProto.AuthorList{
		Authors: parsedAuthors,
	}
}
