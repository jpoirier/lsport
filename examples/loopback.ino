/*
Loop back for testing
 
*/

byte ch = '\0';

void setup()  {
  Serial.begin(115200);
  Serial.println("Starting serial messages:");
  delay(250);
}

void loop() {

  if (Serial.available()) {
	Serial.print("Received: ");
	while(Serial.available()) {
		ch = (byte)Serial.read();
    	Serial.write(ch);
	}
	// Properly terminate the message
	Serial.write('\n');
	Serial.write('\0');

	}
	delay(100);
}

