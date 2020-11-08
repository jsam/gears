package datagears

import (
	"github.com/go-redis/redis/v8"
	"net/url"
	"strconv"
)

//GearsInstance
type GearsInstance struct {
	Addr     string
	Password string
	DB       int

	client *redis.Client
}

//SetDB
func (gears *GearsInstance) SetDB(db int) *GearsInstance {
	gears.DB = db
	return gears
}

//SetPassword
func (gears *GearsInstance) SetPassword(passwd string) *GearsInstance {
	gears.Password = passwd
	return gears
}

//Build
func (gears *GearsInstance) Build() *GearsInstance {
	gears.client = redis.NewClient(&redis.Options{
		Addr:     gears.Addr,
		Password: gears.Password,
		DB:       gears.DB,
	})

	return gears
}

//NewGearsInstance
func NewGearsInstance(addr string) *GearsInstance {
	return &GearsInstance{Addr: addr}
}

func FromRedisDSN(redisDSN string) (*GearsInstance, error) {
	redisUrl, err := url.Parse(redisDSN)
	if err != nil {
		return nil, err
	}

	db, err := strconv.Atoi(redisUrl.Path[1:])
	if err != nil {
		return nil, err
	}

	gearsInstance := NewGearsInstance(redisUrl.Host).
		SetDB(db).
		SetPassword(redisUrl.User.String())

	return gearsInstance, nil
}
