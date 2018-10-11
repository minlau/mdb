package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
)

type DatabaseStore struct {
	databases map[DatabaseGroup]DatabaseInstance
}

func NewDatabaseStore() *DatabaseStore {
	return &DatabaseStore{make(map[DatabaseGroup]DatabaseInstance)}
}

//does not remove added databases if error occurs
func (s *DatabaseStore) AddDatabases(databases []DatabaseConfig) error {
	for _, item := range databases {
		err := s.AddDatabase(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DatabaseStore) AddDatabase(config DatabaseConfig) error {
	if _, ok := s.databases[config.DatabaseGroup]; ok {
		return errors.New("database with given groupId and groupType is already added")
	}

	db, err := openDatabase(config.DatabaseConnConfig)
	if err != nil {
		return errors.Errorf("failed to open database. %v", config)
	}
	s.databases[config.DatabaseGroup] = DatabaseInstance{config, db}
	return nil
}

func (s *DatabaseStore) queryDatabase(groupId int, groupType string, query string) ([]map[string]interface{}, error) {
	if databaseInstance, ok := s.databases[DatabaseGroup{groupId, groupType}]; ok {
		return queryToMap(databaseInstance.DB, query)
	} else {
		return nil, errors.Errorf("no database registered with groupId: %d, groupType: %s", groupId, groupType)
	}
}

//does not have timeout, might be a problem
func (s *DatabaseStore) queryMultipleDatabases(groupType string, query string) []GroupQueryResult {
	var results []GroupQueryResult
	var mutex = &sync.Mutex{}
	var filteredDatabases = make(map[int]*sqlx.DB)

	for key, value := range s.databases {
		if key.GroupType == groupType {
			filteredDatabases[key.GroupId] = value.DB
		}
	}

	c := make(chan int, len(filteredDatabases))

	for groupId, db := range filteredDatabases {
		go func(groupId int, db *sqlx.DB) {
			result, err := queryToMap(db, query)
			var groupQueryResult = GroupQueryResult{
				groupId,
				result,
				err,
			}
			mutex.Lock()
			results = append(results, groupQueryResult)
			mutex.Unlock()
			c <- groupId
		}(groupId, db)
	}
	for range filteredDatabases {
		<-c
	}
	return results
}

func (s *DatabaseStore) getDatabaseItems() []DatabaseDescription {
	arr := make([]DatabaseDescription, 0, len(s.databases))
	for _, value := range s.databases {
		arr = append(arr, value.Config.DatabaseDescription)
	}
	return arr
}

func openDatabase(config DatabaseConnConfig) (*sqlx.DB, error) {
	driverName, err := config.getDriverName()
	if err != nil {
		return nil, err
	}

	connectionUrl, err := config.getConnectionUrl()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driverName, connectionUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return db, nil
}

func queryToMap(db *sqlx.DB, query string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	rows, err := db.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msgf("failed to execute query")
		return nil, err
	}
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			log.Error().Err(err).Msgf("failed to scan data")
			return nil, err
			//continue?
		}

		data = append(data, results)
	}
	return data, nil
}