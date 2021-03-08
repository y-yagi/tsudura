package main

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/y-yagi/configure"
	"github.com/y-yagi/goext/osext"
)

func setupConfig() error {
	var err error

	validateNotExist := func(input string) error {
		if len(input) < 1 {
			return errors.New("Please input value")
		}

		result, err := osext.IsEmptyDir(input)
		if err != nil || !result {
			return errors.New("Please specify empty directory")
		}

		return nil
	}

	validateNotEmpty := func(input string) error {
		if len(input) < 1 {
			return errors.New("Please input value")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Backup directory",
		Validate: validateNotExist,
	}
	if cfg.Root, err = prompt.Run(); err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:    "Bucket",
		Validate: validateNotEmpty,
	}
	if cfg.Bucket, err = prompt.Run(); err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:    "Region",
		Default:  "us-east-1",
		Validate: validateNotEmpty,
	}
	if cfg.Region, err = prompt.Run(); err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:    "Secret",
		Validate: validateNotEmpty,
	}
	if cfg.Secret, err = prompt.Run(); err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:    "Token",
		Validate: validateNotEmpty,
	}
	if cfg.Token, err = prompt.Run(); err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	cfg.Endpoint = "https://s3.wasabisys.com"

	// TODO(y-yagi): validate setting.
	configure.Save(app, &cfg)

	return nil
}
