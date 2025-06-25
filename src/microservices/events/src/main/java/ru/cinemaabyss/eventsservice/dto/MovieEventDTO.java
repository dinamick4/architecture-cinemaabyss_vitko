package ru.cinemaabyss.eventsservice.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class MovieEventDTO extends EventDTO  implements Serializable {

    private Integer movie_id;
    private String title;
    private String action;
}