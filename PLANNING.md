# Ctrl + Shift + Rep
Plan and track workouts with minimal interaction

>  [Theme & Audience](#theme--audience)
>
>  [Planning](#planning)
>> [Form & Content](#form--content)
>> 
>> [Model / View / Controller](#model--view--controller)
>>
>> [Time](#time)
>>
>
> [Stretch goals](#stretch-goals)

## Theme & Audience
I want something that I can mindlessly use at the gym, because when I'm at the gym, I am become mindless. More specifically, I want something that **1)** **plans my workouts**, and **2)** tracks my **weight**, **sets/reps**, and **times how long I spend on each exercise**. The former set of features are to be implemented for the most general of cases, to provide a flexible abstraction for the latter feature. The latter feature will use JSON descriptions of a daily exercise regimen to render the UI.

The audience is just me, but this is contingent on others wanting to use the app.

## Planning
Here, I detail form & content of an exercise regimen, the model/view/controller, and how time is an important consideration throughout

### Form & Content
The first step is to create a JSON schema that describes the form of a "daily exercise regimen". You definitely know whata daily exercise regimen is, but in case you don't, please refer to the following

```
Me > What is an exercise regimen?

ChatGPT > An exercise regimen is a structured plan or schedule of physical activities designed to improve or maintain fitness, health, or performance. It typically includes details like:

- Type of exercises (e.g., cardio, strength training, flexibility)
- Frequency (how many days per week)
- Duration (how long each session lasts)
- Intensity (how hard you work—light, moderate, vigorous)
- Progression (how the plan evolves over time)

[...]
```

Chat described a lot of things, but I care about the **type of exercise**, **sets**, and **reps**. The other things are either not applicable or are discussed later. 

A JSON document describing such a plan might look like this.\
**NOTE**: The null value means *"Do as many reps as you can"*

```Back and Bicep Day.json```
```json
{ 
    "back":
    [
        {
            "exercise" : "Deadlifts",
            "sets"     : 4,
            "reps"     : 5
        },
        {
            "exercise" : "Lat Pulldowns",
            "sets"     : 3,
            "reps"     : null
        },
        {
            "exercise" : "Bent-Over Barbell Rows",
            "sets"     : 3,
            "reps"     : 10
        },
        {
            "exercise" : "Single-Arm Dumbbell Rows",
            "sets"     : 3,
            "reps"     : 10
        },
        {
            "exercise" : "Rear Delt Flyes",
            "sets"     : 3,
            "reps"     : 10
        }
    ],
    "biceps" :
    [
        {
            "exercise" : "Barbell Curls",
            "sets"     : 3,
            "reps"     : 10
        },
        {
            "exercise" : "Incline Dumbbell Curls",
            "sets"     : 3,
            "reps"     : 10
        },
        {
            "exercise" : "Cable Hammer Curls",
            "sets"     : 3,
            "reps"     : 10
        }
    ]
}
```

Making a schema for this data format is the next step. If you start by writing a schema, you will reap many rewards including vscode autofill, robust validatation, and useful insights into how you will need to structure data in the application. But the best part about making a schema is that you can write tests based on the rules in the schema. The tests essentially write themselves. 

### Model / View / Controller
At the core of the data model is a list of objects. There exists one object for each calendar day since the user started using with the application. Each object contains an exercise regimen object, and a clone of that exercise regimen. The clone contains additional properties that are collected throughout the exercise. These additional properties are detailed in the controller section

The view will be very simple. Each view will persistent an identical navigation element ad the head of the page. The navigation bar will have two elements. The first is for viewing historical exercise data, and this will be as mundane as you'd imagine. The second is for interacting with the app during an exercise. The views change throughout the exercise, updating as the user proceeds, and guiding the user through the regimen. Future additions include ~~encouraging / bullying the user based on how they are performing~~

The controller is very minimal. The view presents the current exercise, and some number inputs in a form. These are for the weight/reps/sets. A big green **NEXT** button is used to advance the view to the next exercise. I don't want to spend my rest period doing manual data entry, so the weight/reps/sets are populated with from the previous instance of this specific exercise regimen. 99% of the time, I can just click **NEXT** through my exercise. One button. One braincell. It sells itself. 

## Time 
Each time the user clicks the next button, a timestamp is recorded. These timesteps, if properly recorded, can be used to derive many useful values, some of which I have not enumerated here because I still haven't thought of them

```
 1. Total time spent resting 
 2. Total time spent exercising 

 3. Time spent resting for a specific exercise
 4. Time spent exercising for a specific exercise

 5. Statistics of 1 and 2
 6. Statistics of 3 and 4
 ```

## Stretch Goals

### Graphical view 
It would be nice to present statistics graphically. Maybe a bar chart? Maybe a graph? Maybe something else, but ***definitely a strategy pattern!*** 

### Friends 
My first idea was color coding different lines foreach friend, and compositing the result onto the SVG view. Then I remembered that sharing data between users is fraught with traps. For instance, people get competitive about quantifiable statistics. I don't know why this is a thing, but I will not have my code enable this behavior.Competition might discourage the "losers", and this is an unacceptable outcome. Ethics aside, each sample was sampled from a unique source. Unique sources produce incommensurate data. Eo ipso, the data would be meaningless.

My second idea is was, if a user has been consistently exceeding their baseline, their friends will get a notification, instructing them to encourage the user. This change would replace meaningless numbers and meaningless competitions with a collaborative social network. When you're doing well, everybody will know. If not, radio silence. 
