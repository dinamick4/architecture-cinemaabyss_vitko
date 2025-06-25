package ru.cinemaabyss.eventsservice.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import ru.cinemaabyss.eventsservice.dto.EventDTO;

@Service
public class EventProducerService {

    @Autowired
    private KafkaTemplate<String, EventDTO> kafkaTemplate;

    public void send(String topic, EventDTO event) {
        kafkaTemplate.send(topic, event);
    }
}