package ru.yandex.practicum.smarthome.components;

import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;
import ru.yandex.practicum.smarthome.dto.CommandDto;

@Component
public class KafkaProducer {

    private final KafkaTemplate<String, CommandDto> kafkaTemplate;
    private final String TOPIC_NAME = "commands";

    public KafkaProducer(KafkaTemplate<String, CommandDto> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendMessage(CommandDto message) {
        kafkaTemplate.send(TOPIC_NAME, message);
        System.out.println("Message has been sucessfully sent to the topic: " + TOPIC_NAME);
    }
}