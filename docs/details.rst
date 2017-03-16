======================
Details about Sisyphus
======================

In a nutshell this is what sisyphus is doing after starting it:

1. Load all mails in the good and junk directories
2. Check for each mail whether it has been processed
3. If not processed, than classify it based on subject and body content
4. If processed, check whether it has not been moved from junk to good or vice versa
5. If a processed mail has been moved to the other folder, unlearn its words from the old one and learn them for the new classification
6. Observe the junk and good mail directory in real time and handle new incoming mails directly
