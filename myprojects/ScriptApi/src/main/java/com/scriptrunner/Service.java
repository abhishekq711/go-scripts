package com.scriptrunner;

import org.springframework.stereotype.Component;
import com.scriptrunner.AttributeModel;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

@Component
public class Service {

    public static void main(String[] args){}

    public void run(AttributeModel attributes) {
        String[] cmd = new String[] {"/home/atuljadhav/Desktop/GenerateCode/run.sh","env="+attributes.getEnv(),"ide="+attributes.getIde(),"identity="+attributes.getIdentity()};
        ProcessBuilder pb = new ProcessBuilder(cmd);
        try {
            Process p = pb.start();
            BufferedReader reader = new BufferedReader(new InputStreamReader(p.getInputStream()));
            String s = null;
            while ((s = reader.readLine()) != null){
                System.out.println(s);
            }
        } catch (IOException e){
            e.printStackTrace();
            e.getMessage();
        }
    }
}
