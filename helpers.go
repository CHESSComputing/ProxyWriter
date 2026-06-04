package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CHESSComputing/golib/beamlines"
	srvConfig "github.com/CHESSComputing/golib/config"
)

func mergeData(userRecord map[string]any) map[string]any {
	rec := make(map[string]any)
	// TODO:
	// 1. for given user data locate associated main record on /nfs using record keys: btr+sample_name
	// the CHESS /nfs path should be constructed with some rule
	// 2. merge two records together
	return rec
}

func validateData(rec map[string]any) error {
	return nil
}

func insertData(userRecord map[string]any) (string, error) {

	// merge user record with spec record
	rec := mergeData(userRecord)

	var sname string
	if val, ok := rec["SchemaName"]; ok {
		sname = fmt.Sprintf("%s", val)
	} else {
		return "", errors.New("provided record does not contain schema name")
	}

	// look up schema file name from schema file
	var schemaFile string
	for _, sfile := range srvConfig.Config.CHESSMetaData.SchemaFiles {
		if strings.Contains(strings.ToLower(sname), strings.ToLower(sfile)) {
			schemaFile = sfile
		}
	}

	if schemaFile == "" {
		msg := fmt.Sprintf("for schema name %s no schema file found", sname)
		return "", errors.New(msg)
	}

	// validate merged record
	schema := beamlines.Schema{FileName: schemaFile}
	err := schema.Load()
	if report := schema.ValidateAll(rec); report != "" {
		return "", errors.New(report)
	}

	var did string
	if val, ok := rec["did"]; ok {
		did = fmt.Sprintf("%v", val)
	} else {
		msg := "provided metadata record does not contain did attribute"
		return "", errors.New(msg)
	}

	// insert record to metaDB
	err = metaDB.InsertRecord(
		srvConfig.Config.UserMetaData.MongoDB.DBName,
		srvConfig.Config.UserMetaData.MongoDB.DBColl,
		rec)

	return did, err
}
