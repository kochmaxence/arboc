package arboc

import (
	"reflect"
	"unsafe"

	"github.com/fatih/structtag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type fieldConfig struct {
	Long        string
	Short       string
	Description string
	Persistent  bool
	Required    bool
}

func getFieldConfig(field *fieldData, tags *structtag.Tags) (*fieldConfig, bool) {
	config := &fieldConfig{}

	cmdTag, _ := tags.Get("cmd")
	if cmdTag == nil {
		return nil, true
	}

	config.Long = cmdTag.Name
	argsCount := len(cmdTag.Options)

	if argsCount > 0 {
		config.Short = cmdTag.Options[0]

	}

	descTag, _ := tags.Get("desc")
	if descTag != nil {
		config.Description = descTag.Name
	}

	persistTag, _ := tags.Get("persist")
	if persistTag != nil {
		val := persistTag.Name

		if val == "true" || val == "1" {
			config.Persistent = true
		}
	}

	requiredTag, _ := tags.Get("required")
	if requiredTag != nil {
		val := requiredTag.Name

		if val == "true" || val == "1" {
			config.Required = true
		}
	}

	return config, false
}

func getFlagSet(cmd *cobra.Command, cfg *fieldConfig) *pflag.FlagSet {
	if cfg.Persistent {
		return cmd.PersistentFlags()
	} else {
		return cmd.Flags()
	}
}

func handleMarkRequired(cmd *cobra.Command, cfg *fieldConfig) {
	if cfg.Required {
		if cfg.Persistent {
			cmd.MarkPersistentFlagRequired(cfg.Long)
		} else {
			cmd.MarkFlagRequired(cfg.Long)
		}
	}
}

func newDefaultFlagger() *Flagger {
	funcs := make(map[reflect.Kind]flaggerFunc, 0)

	funcs[reflect.Bool] = func(cmd *cobra.Command, field *fieldData, tags *structtag.Tags) {
		cfg, skip := getFieldConfig(field, tags)

		if skip {
			return
		}

		ptr := (*bool)(unsafe.Pointer(field.Value.Addr().Pointer()))
		flags := getFlagSet(cmd, cfg)

		if cfg.Short != "" {
			flags.BoolVarP(ptr, cfg.Long, cfg.Short, field.Value.Bool(), cfg.Description)
		} else {
			flags.BoolVar(ptr, cfg.Long, field.Value.Bool(), cfg.Description)
		}

		handleMarkRequired(cmd, cfg)
	}

	funcs[reflect.Int] = func(cmd *cobra.Command, field *fieldData, tags *structtag.Tags) {
		cfg, skip := getFieldConfig(field, tags)

		if skip {
			return
		}

		ptr := (*int)(unsafe.Pointer(field.Value.Addr().Pointer()))
		flags := getFlagSet(cmd, cfg)

		if cfg.Short != "" {
			flags.IntVarP(ptr, cfg.Long, cfg.Short, int(field.Value.Int()), cfg.Description)
		} else {
			flags.IntVar(ptr, cfg.Long, int(field.Value.Int()), cfg.Description)
		}

		handleMarkRequired(cmd, cfg)
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = funcs[reflect.Int]

	funcs[reflect.String] = func(cmd *cobra.Command, field *fieldData, tags *structtag.Tags) {
		cfg, skip := getFieldConfig(field, tags)

		if skip {
			return
		}

		ptr := (*string)(unsafe.Pointer(field.Value.Addr().Pointer()))
		flags := getFlagSet(cmd, cfg)

		if cfg.Short != "" {
			flags.StringVarP(ptr, cfg.Long, cfg.Short, field.Value.String(), cfg.Description)
		} else {
			flags.StringVar(ptr, cfg.Long, field.Value.String(), cfg.Description)
		}

		handleMarkRequired(cmd, cfg)
	}

	funcs[reflect.Slice] = func(cmd *cobra.Command, field *fieldData, tags *structtag.Tags) {
		cfg, skip := getFieldConfig(field, tags)

		if skip {
			return
		}

		flags := getFlagSet(cmd, cfg)

		switch field.Value.Type().Elem().Kind() {
		case reflect.String:
			ptr := (*[]string)(unsafe.Pointer(field.Value.Addr().Pointer()))

			values := field.Value.Interface().([]string)

			if cfg.Short != "" {
				flags.StringSliceVarP(ptr, cfg.Long, cfg.Short, values, cfg.Description)
			} else {
				flags.StringSliceVar(ptr, cfg.Long, values, cfg.Description)
			}
		}

		handleMarkRequired(cmd, cfg)
	}

	return &Flagger{
		FuncByKind: funcs,
	}
}
