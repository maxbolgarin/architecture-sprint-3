package ru.yandex.practicum.smarthome.components;

import org.springframework.stereotype.Component;
import ru.yandex.practicum.smarthome.dto.CommandDto;

@Component
public class DeviceType {
    // TODO: get these args from device types microservice
    private final String heatingSystemOnId = "heating_system_on";
    private final String heatingSystemOffId = "heating_system_off";
    private final String heatingSystemSetId = "heating_system_set";
    private final String plainCommandType = "plain";

    public CommandDto GetHeatingSystemOnCommand() {
        return new CommandDto(heatingSystemOnId, plainCommandType, "");
    }
    public CommandDto GetHeatingSystemOffCommand() {
        return new CommandDto(heatingSystemOffId, plainCommandType, "");
    }
    public CommandDto GetHeatingSystemSetCommand(double temperature) {
        return new CommandDto(heatingSystemSetId, plainCommandType, String.valueOf(temperature));
    }   
}