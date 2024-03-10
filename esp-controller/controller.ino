#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>
#include <PubSubClient.h>

#define LED D2 // Red Led

// Wireless connectivity 
const char* ssid = "FRITZ!Box 7530 QZ";
const char* password = "65367332125656334037";
WiFiClient espClient;
PubSubClient client(espClient);

// MQTT Broker ip and port details, set to loopback as an example 
const char *mqtt_broker = "127.0.0.1"; 
const int mqtt_port = 1883;
// MQTT topic used to publish sensor readings 
const char *topic = "pir_sensor/motion_reading";
// microcontroller id for MQTT
const char *message = "movement_detected-sensor2"; 

// PIR sensor 
int sensor = 13;  // Digital pin D7

// ultrasonic trigger pin and echo pin 
const int trigPin = D4;
const int echoPin = D8;

// duration and distance variable 
long duration = 0; 
int distance = 0; 

// set the port for webserver
ESP8266WebServer server(80);

// the setup function runs once when you press reset or power the board
void setup() {
  Serial.begin(115200);

  // start by connecting to a wireless network
  Serial.print("I am auxiliary sensor. I will attempt to connect to wireless network");
  Serial.println(ssid);
  
  // set the digital pin as output
  pinMode(LED, OUTPUT);

  WiFi.begin(ssid, password);
  
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }

  Serial.println("");
  Serial.println("WiFi connected");  
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());

  // set the trigger pin as output and echo pin as input
  pinMode(trigPin, OUTPUT);
  pinMode(echoPin, INPUT);

  // declare PIR sensor as input
  pinMode(sensor, INPUT);   

  //set the webserver root page 
  server.on("/", get_index);

  //start the server
  server.begin();
  Serial.println("Server listening!");

  // connecting to a mqtt broker
  client.setServer(mqtt_broker, mqtt_port);
  client.setCallback(callback);

  while (!client.connected())
  {
      String client_id = "esp8266-client-";
      client_id += String(WiFi.macAddress());

      Serial.printf("The client %s connects to mosquitto mqtt broker\n", client_id.c_str());

      if (client.connect(client_id.c_str()))
      {
          Serial.println("Public emqx mqtt broker connected");
      }
      else
      {
          Serial.print("failed with state ");
          Serial.print(client.state());
          delay(2000);
      }
  }

  // publish Hello
  client.publish(topic, "Hello from Auxiliary controller");

}

void callback(char *topic, byte *payload, unsigned int length)
{

    for (int i = 0; i < length; i++)
    {
        Serial.print((char)payload[i]);
    }

    Serial.println();
    Serial.println(" - - - - - - - - - - - -");
}

// the loop function runs over and over again forever
void loop() {

  server.handleClient();  
  client.loop();

  delay(500);

  // main task - deetect movement 
  detectMovement();

}

void detectMovement(){
   
  int previous_distance = distanceCentimeter();
  
  Serial.print("previous distance: ");
  Serial.print(previous_distance);
  Serial.println(" centimeters");

  // slow down readings  
  delay(1000);
  int new_distance = distanceCentimeter();

  Serial.print("new distance: ");
  Serial.print(new_distance);
  Serial.println(" centimeters");

  if (abs(new_distance - previous_distance) > 2)
  {
    Serial.println("");
    Serial.println("------movement detected (details below)------");
    String distanceStr = String(distance);
    Serial.print("previous distance: ");
    Serial.print(previous_distance);
    Serial.print(";new distance: ");
    Serial.print(new_distance);
    Serial.println("");
    Serial.println("------movement detected, let's publish------");

    client.publish(topic, message);

    // let's sleep for 20 seconds 
    delay(20000);
  }

}

int distanceCentimeter(){

  //clear the trigPin
  digitalWrite(trigPin, LOW);
  delayMicroseconds(2);

  //sets the trigPin on HIGH state for 10micro seconds
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);

  // clears the trigPin
  digitalWrite(trigPin, LOW);

  // read from the echo pin , returns the sound wave travel time which has to be converted 
  duration = pulseIn(echoPin, HIGH);

  // distance calculation 
  distance =  (duration * 0.034) /2 ;

  return distance;
      
}

// basic http service, useful for debugging 
void get_index(){
  String html ="<!DOCTYPE html> <html> ";
  html += "<head><meta http-equiv=\"refresh\" content=\"2\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"></head>";
  html += "<body> <h1>Auxiliary controller</h1>";
  html +="<p> Welcome to Auxiliary controller dashboard </p>";
  html += "<div> <p>  If you can see the dashboard and the following reading it means the controller is functioning correctly: <strong>";
  html += distance;
  html +="</strong> cm. </p> </div>";
  html +="</body> </html>";

  server.send(200, "text/html", html);

}