package com.scriptrunner;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/script")
public class RunController {

    @Autowired
    private Service runner;

    @PostMapping(value = "/run")
    public void runScript(@RequestBody AttributeModel attributeModel){
            runner.run(attributeModel);
    }

}
