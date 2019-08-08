# mkhousecurve

This utility takes a text frequency response exported from Room EQ Wizard and
converts it into a house curve for use in generating EQ settings with Room EQ
Wizard.

For example:

```bash
mkhousecurve -delim " " -in vemonkplus.txt -out vemonkplus_housecurve.txt -comment "House curve based on VE Monk Plus"
```

The frequency response can be exported to text from the menu File -> Export -> Export Measurement as Text.
This program assumes that the frequency response is tab delimited.
