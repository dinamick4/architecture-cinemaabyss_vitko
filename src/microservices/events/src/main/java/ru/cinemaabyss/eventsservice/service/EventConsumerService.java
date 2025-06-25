package ru.cinemaabyss.eventsservice.service;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
public class EventConsumerService {

    private static final Logger logger = LoggerFactory.getLogger(EventConsumerService.class);

    @KafkaListener(topics = {"user-events", "payment-events", "movie-events"})
    public void consume(ConsumerRecord<String, String> record){
        logger.info("Received message from Kafka: topic={}, partition={}, offset={}",
                record.topic(), record.partition(), record.offset());
        logger.info("Value: {}", record.value());
    }
}