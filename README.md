See a demo and a walkthrough of the code here: https://www.youtube.com/watch?v=5ONrTsfzzu0&t=42s

I like to have high and low tide on my calendar, so I created this simple app to scrape the tide data for the year from tidetime.org and then populate an iCal file with the data.

The app can create three different calendars: high tide only, low tide only, or both high and low tide. 

You can choose custom colors for both high and low tide, but some calendar applications, such as Google Calendar, will not respect the custom event colors, so if you want different colors for high and low tide, it is best to import the two .ics files individually and assign a different color to each.

At the moment, to change which calendar is generated, you need to edit the source code by modifying the `whichTides` variable found at the bottom fo the `main` function, and you can change the custom colors by modifying the `color` map found at the top of the `createCalendar` function.

I might make this more usable by adding either command-line arguments or a web interface for selecting calendar options... maybe if I expand the app to allow the user to select the region, then it will make sense to improve the UX for the general public.

Have a great day!
