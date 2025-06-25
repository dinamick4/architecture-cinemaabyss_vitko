package ru.cinemaabyss.eventsservice.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import ru.cinemaabyss.eventsservice.dto.MovieEventDTO;
import ru.cinemaabyss.eventsservice.dto.PaymentEventDTO;
import ru.cinemaabyss.eventsservice.dto.UserEventDTO;
import ru.cinemaabyss.eventsservice.service.EventProducerService;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/events")
public class EventsController {

    @Autowired
    private EventProducerService eventProducerService;

    @GetMapping("/health")
    public ResponseEntity<Map<String, Boolean>> health() {
        Map<String, Boolean> response = new HashMap<>();
        response.put("status", true);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/user")
    public ResponseEntity<Map<String, String>> createUserEvent(@RequestBody UserEventDTO eventDto) {
        eventProducerService.send("user-events", eventDto);
        Map<String, String> response = new HashMap<>();
        response.put("status", "success");
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @PostMapping("/payment")
    public ResponseEntity<Map<String, String>> createPaymentEvent(@RequestBody PaymentEventDTO eventDto) {
        eventProducerService.send("payment-events", eventDto);
        Map<String, String> response = new HashMap<>();
        response.put("status", "success");
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @PostMapping("/movie")
    public ResponseEntity<Map<String, String>> createMovieEvent(@RequestBody MovieEventDTO eventDto) {
        eventProducerService.send("movie-events", eventDto);
        Map<String, String> response = new HashMap<>();
        response.put("status", "success");
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }
}