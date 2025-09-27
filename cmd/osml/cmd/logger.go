package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/logger"
	"github.com/willie68/osmltools/internal/sdformatter"
)

var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "commands to work with the osmlogger configuration",
	Long:  `commands to work with the open sea map logger configuration`,
}

var loggerReadCmd = &cobra.Command{
	Use:   "read",
	Short: "read the osmlogger configuration",
	Long:  `read the open sea map logger configuration`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := logger.ReadFromSDCard(sdCardFolder)
		if err != nil {
			return err
		}
		if JSONOutput {
			js, err := cfg.JSON()
			if err != nil {
				return err
			}
			fmt.Println(js)
			return nil
		}
		fmt.Printf("configuration read from %s/config.dat\n", sdCardFolder)
		fmt.Printf("config: %s\n", cfg.String())
		return nil
	},
}

var loggerWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "write the osmlogger configuration",
	Long:  `write the open sea map logger configuration`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		seatalk, _ := cmd.Flags().GetBool("seatalk")
		baudA, _ := cmd.Flags().GetInt16("baudA")
		baudB, _ := cmd.Flags().GetInt16("baudB")
		vesselID, _ := cmd.Flags().GetInt16("vesselid")
		gyro, _ := cmd.Flags().GetBool("gyro")
		supply, _ := cmd.Flags().GetBool("supply")
		cfg := logger.NewLoggerConfig().
			WithBaudA(baudA).
			WithBaudB(baudB).
			WithGyro(gyro).
			WithSupply(supply).
			WithSeatalk(seatalk).
			WithVesselID(vesselID)
		err := cfg.Validate()
		if err != nil {
			return err
		}
		sdformat, _ := cmd.Flags().GetBool("sdformat")
		if sdformat {
			OutputWithJSONCheckf("formatting sd card at %s\n", sdCardFolder)
			err = sdformatter.FormatFAT32(sdCardFolder)
			if err != nil {
				return err
			}
			OutputWithJSONCheckf("sd card at %s formatted\n", sdCardFolder)
		}
		sdlabel, _ := cmd.Flags().GetString("sdlabel")
		if sdlabel != "" {
			OutputWithJSONCheckf("setting the label of sd card to %s on %s\n", sdlabel, sdCardFolder)
			err = sdformatter.SetLabel(sdCardFolder, sdlabel)
			if err != nil {
				return err
			}
			OutputWithJSONCheckf("label set to %s on %s\n", sdlabel, sdCardFolder)
		}
		err = cfg.WriteToSDCard(sdCardFolder)
		if err != nil {
			return err
		}
		if JSONOutput {
			js, err := cfg.JSON()
			if err != nil {
				return err
			}
			fmt.Println(js)
			return nil
		}
		OutputWithJSONCheckf("configuration written to %s/config.dat\n", sdCardFolder)
		OutputWithJSONCheckf("config: %s\n", cfg.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loggerCmd)

	loggerCmd.AddCommand(loggerReadCmd)

	loggerCmd.AddCommand(loggerWriteCmd)

	loggerWriteCmd.Flags().BoolP("seatalk", "", false, "channel A accepts seatalk input")
	loggerWriteCmd.Flags().Int16P("baudA", "a", 4800, "channel A communication baud rate")
	loggerWriteCmd.Flags().Int16P("baudB", "b", 4800, "channel B communication baud rate")
	loggerWriteCmd.Flags().Int16P("vesselid", "", 0, "id of the vessel to set or get the configuration")
	loggerWriteCmd.Flags().BoolP("gyro", "", true, "write internal gyro data to the data files")
	loggerWriteCmd.Flags().BoolP("supply", "", false, "write internal supply data to the data files")
	loggerWriteCmd.Flags().Bool("sdformat", false, "format the sd card before writing the configuration")
	loggerWriteCmd.Flags().String("sdlabel", "", "setting the label of the sd card")
}
