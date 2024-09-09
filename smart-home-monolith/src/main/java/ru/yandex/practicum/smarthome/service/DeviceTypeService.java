package ru.yandex.practicum.smarthome.service;

import org.springframework.stereotype.Service;
import ru.yandex.practicum.smarthome.dto.CommandDto;

@Service
public class DeviceTypeService {
    // TODO: get these args from device types microservice
    private String heatingSystemOnId = "heating_system_on";
    private String heatingSystemOffId = "heating_system_off";
    private String heatingSystemSetId = "heating_system_set";
    private String plainCommandType = "plain";

    public CommandDto GetHeatingSystemOnCommand() {
        return new CommandDto(this.heatingSystemOnId, this.plainCommandType, "");
    }
    public CommandDto GetHeatingSystemOffCommand() {
        return new CommandDto(this.heatingSystemOffId, this.plainCommandType, "");
    }
    public CommandDto GetHeatingSystemSetCommand(double temperature) {
        return new CommandDto(this.heatingSystemSetId, this.plainCommandType, String.valueOf(temperature));
    }   
}