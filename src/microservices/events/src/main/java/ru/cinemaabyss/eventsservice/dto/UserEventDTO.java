package ru.cinemaabyss.eventsservice.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class UserEventDTO extends EventDTO implements Serializable {

    private String username;
    private String action;
    private String timestamp;
}