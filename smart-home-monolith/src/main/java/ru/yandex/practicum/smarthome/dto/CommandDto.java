package ru.yandex.practicum.smarthome.dto;

import lombok.Data;
import java.util.Random;
import java.util.Date;
import java.util.UUID;

@Data
public class CommandDto {
    private String commandId;
    private String commandTypeId;
    private Date createTime;
    private String code;
    private String commandType;

    public CommandDto(String typeId, String commandType, String code) {
        this.commandId = UUID.randomUUID().toString().replaceAll("_", "");
        this.commandTypeId = typeId;
        this.createTime = new Date();
        this.commandType = commandType;
        this.code = code;
    }
}
