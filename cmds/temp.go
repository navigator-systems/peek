package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shirou/gopsutil/v4/sensors"
)

type temperatureInfo struct {
	SensorKey   string  `json:"sensor_key"`
	Temperature float64 `json:"temperature"`
	High        float64 `json:"high"`
	Critical    float64 `json:"critical"`
}

func TEMP(jsonOutput bool) error {
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		return fmt.Errorf("unable to get temperature info: %w", err)
	}

	info := make([]temperatureInfo, 0, len(temperatures))
	for _, temperature := range temperatures {
		info = append(info, temperatureInfo{
			SensorKey:   temperature.SensorKey,
			Temperature: temperature.Temperature,
			High:        temperature.High,
			Critical:    temperature.Critical,
		})
	}

	if jsonOutput {
		encoded, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to encode temperature info as JSON: %w", err)
		}
		fmt.Println(string(encoded))
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SENSOR\tTEMP(C)\tHIGH(C)\tCRITICAL(C)")
	for _, temperature := range info {
		fmt.Fprintf(w, "%s\t%.1f\t%.1f\t%.1f\n",
			temperature.SensorKey,
			temperature.Temperature,
			temperature.High,
			temperature.Critical,
		)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("unable to print temperature table: %w", err)
	}
	return nil
}
