package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

// https://discord.com/developers/docs/reference#snowflakes
const discordEpoch int64 = 1420070400000

type snowflakeError struct {
	snowflake int64
	message   string
}

func (e *snowflakeError) Error() string {
	return fmt.Sprintf("Snowflake ID %d is invalid: %s", e.snowflake, e.message)
}

func newSnowflakeError(snowflake int64, message string) error {
	return &snowflakeError{snowflake: snowflake, message: message}
}

// converts a Discord snowflake ID to a Unix timestamp (in milliseconds) using the provided epoch
func convertSnowflakeToDate(snowflake int64, epoch int64) int64 {
	milliseconds := snowflake >> 22
	return milliseconds + epoch
}

// validates a Discord snowflake ID and returns a Unix timestamp (in milliseconds) if valid
func validateSnowflake(snowflake int64, epoch int64) (int64, error) {
	if snowflake < 0 {
		return 0, newSnowflakeError(snowflake, "snowflake IDs must be positive integers")
	}

	if snowflake < 4194304 {
		return 0, newSnowflakeError(snowflake, "snowflake is too small")
	}

	timestamp := convertSnowflakeToDate(snowflake, epoch)
	if timestamp < 0 {
		return 0, newSnowflakeError(snowflake, "snowflake has fewer digits than expected")
	}

	return timestamp, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a Discord snowflake ID: ")
		snowflake, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		snowflake = snowflake[:len(snowflake)-2] // remove \n and \r
		intSnowflake, err := strconv.ParseInt(snowflake, 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		timestamp, err := validateSnowflake(intSnowflake, discordEpoch)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dateTime := time.UnixMilli(timestamp)
		fmt.Println(dateTime.Format(time.RFC1123))
		break
	}
}
