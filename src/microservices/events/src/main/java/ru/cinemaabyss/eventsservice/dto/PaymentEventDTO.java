package ru.cinemaabyss.eventsservice.dto;

import lombok.Data;

import java.io.Serializable;
import java.math.BigDecimal;

@Data
public class PaymentEventDTO extends EventDTO implements Serializable {

    private Integer payment_id;
    private BigDecimal amount;
    private String status;
    private String timestamp;
    private String method_type;
}