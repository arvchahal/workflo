# CLI Application Built Using the BubbleTea Framework

This directory contains the code for a CLI application built using the [BubbleTea](https://github.com/charmbracelet/bubbletea) framework. The application manages state transitions and user interactions to provide a seamless command-line experience.

Due to the nature of the BubbleTea framework, the code may be complex or unfamiliar to new contributors. We recommend reviewing the BubbleTea documentation and examples in the official repository to become acquainted with its concepts and patterns.

## Overview

The application is structured around states and transitions, managed through the BubbleTea framework. The main components of the application are:

- **`model.go`**: Defines the application's data model and the various states.
- **`init.go`**: Initializes the model and sets up initial components.
- **`update.go`**: Contains the logic for state transitions and handles user inputs.
- **`view.go`**: Manages how the application renders the user interface based on the current state.

## Adding New Features

To add new features or states to the CLI application, you will need to update the following files:

- `model.go`
- `init.go`
- `update.go`
- `view.go`

Below is a detailed guide on how to update each file to incorporate new features.

---

### 1. `model.go`

This file defines the application's data model, including the possible states and any data structures used throughout the application.

#### Steps to Update `model.go`:

- **Define New States**: Add any new states required by your feature to the `state` type definition.

  ```go
  // model.go

  // Define the application states
  type state int

  const (
      // Existing states
      stateWorkflowName state = iota
      stateRunner
      // Add your new state(s) here
      stateNewFeatureInput
      stateNewFeatureConfirmation
      // ...
  )
  ```

- **Update the Model Struct**: If your feature requires new data fields, add them to the `model` struct.

  ```go
  // model.go

  type model struct {
      // Existing fields
      state       state
      textInput   textinput.Model
      // ...

      // New fields for your feature
      newFeatureInput     textinput.Model
      newFeatureData      string
      newFeatureConfirmed bool
  }
  ```

---

### 2. `init.go`

This file initializes the model and sets up any initial components, such as text inputs or lists.

#### Steps to Update `init.go`:

- **Initialize New Components**: If your feature uses new UI components (e.g., text inputs, lists), initialize them here.

  ```go
  // init.go

  func NewModel() model {
      // Existing initializations
      ti := textinput.New()
      ti.Placeholder = "Enter workflow name"
      // ...

      // Initialize new components
      newFeatureInput := textinput.New()
      newFeatureInput.Placeholder = "Enter data for new feature"
      newFeatureInput.CharLimit = 100
      newFeatureInput.Width = 50

      return model{
          // Existing fields
          state:     stateWorkflowName,
          textInput: ti,
          // ...

          // New fields
          newFeatureInput: newFeatureInput,
      }
  }
  ```

- **Set Default Values**: Initialize any new fields with default values if necessary.

---

### 3. `update.go`

This file contains the main logic for handling user inputs and updating the application's state.

#### Steps to Update `update.go`:

- **Handle New States**: Add cases to the `Update` function to handle your new states.

  ```go
  // update.go

  func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      var cmd tea.Cmd

      switch m.state {
      // Existing states
      case stateWorkflowName:
          // Existing logic
          // ...

      // New feature input state
      case stateNewFeatureInput:
          m.newFeatureInput, cmd = m.newFeatureInput.Update(msg)
          return m.handleNewFeatureInputState(msg, cmd)

      // New feature confirmation state
      case stateNewFeatureConfirmation:
          return m.handleNewFeatureConfirmationState(msg, cmd)

      // ...
      }

      return m, cmd
  }
  ```

- **Implement State Handlers**: Write functions to handle logic specific to your new states.

  ```go
  // update.go

  func (m model) handleNewFeatureInputState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.KeyMsg:
          switch msg.String() {
          case "enter":
              m.newFeatureData = m.newFeatureInput.Value()
              m.newFeatureInput.Reset()
              m.state = stateNewFeatureConfirmation
              return m, textinput.Blink
          case "ctrl+c", "q":
              return m, tea.Quit
          }
      }
      return m, cmd
  }

  func (m model) handleNewFeatureConfirmationState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.KeyMsg:
          switch msg.String() {
          case "y", "Y":
              m.newFeatureConfirmed = true
              // Transition to next state or perform action
              m.state = stateNext
              return m, cmd
          case "n", "N":
              m.newFeatureConfirmed = false
              // Handle negative confirmation
              m.state = stateWorkflowName // Or another appropriate state
              return m, cmd
          case "ctrl+c", "q":
              return m, tea.Quit
          }
      }
      return m, cmd
  }
  ```

- **Update the Main `Update` Function**: Ensure that the new states are properly connected within the state machine.

---

### 4. `view.go`

This file handles the rendering of the user interface based on the current state.

#### Steps to Update `view.go`:

- **Render New States**: Add cases to the `View` function to display the UI for your new states.

  ```go
  // view.go

  func (m model) View() string {
      switch m.state {
      // Existing states
      case stateWorkflowName:
          return fmt.Sprintf("Enter workflow name:\n\n%s\n\n(Press Enter to continue)", m.textInput.View())

      // New feature input state
      case stateNewFeatureInput:
          return fmt.Sprintf("Enter data for new feature:\n\n%s\n\n(Press Enter to continue)", m.newFeatureInput.View())

      // New feature confirmation state
      case stateNewFeatureConfirmation:
          return fmt.Sprintf("You entered: %s\n\nConfirm? (y/n)", m.newFeatureData)

      // ...
      }
  }
  ```

---

## Example: Adding a "Custom Message" Feature

Suppose you want to add a feature where the user can input a custom message that will be displayed later.

### Steps to Add the "Custom Message" Feature:

1. **Update `model.go`**:

   ```go
   // model.go

   // Add new state
   const (
       // Existing states
       // ...
       stateCustomMessageInput state = iota + existingStateCount
       stateDisplayCustomMessage
   )

   type model struct {
       // Existing fields
       // ...

       // New fields
       customMessageInput textinput.Model
       customMessage      string
   }
   ```

2. **Update `init.go`**:

   ```go
   // init.go

   func NewModel() model {
       // Existing initializations
       // ...

       // Initialize new text input for custom message
       customMessageInput := textinput.New()
       customMessageInput.Placeholder = "Enter your custom message"
       customMessageInput.CharLimit = 200
       customMessageInput.Width = 50

       return model{
           // Existing fields
           // ...

           // New fields
           customMessageInput: customMessageInput,
       }
   }
   ```

3. **Update `update.go`**:

   ```go
   // update.go

   func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       var cmd tea.Cmd

       switch m.state {
       // Existing states
       // ...

       case stateCustomMessageInput:
           m.customMessageInput, cmd = m.customMessageInput.Update(msg)
           return m.handleCustomMessageInputState(msg, cmd)

       case stateDisplayCustomMessage:
           return m.handleDisplayCustomMessageState(msg, cmd)
       }

       return m, cmd
   }

   func (m model) handleCustomMessageInputState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
       switch msg := msg.(type) {
       case tea.KeyMsg:
           switch msg.String() {
           case "enter":
               m.customMessage = m.customMessageInput.Value()
               m.customMessageInput.Reset()
               m.state = stateDisplayCustomMessage
               return m, cmd
           case "ctrl+c", "q":
               return m, tea.Quit
           }
       }
       return m, cmd
   }

   func (m model) handleDisplayCustomMessageState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
       switch msg := msg.(type) {
       case tea.KeyMsg:
           switch msg.String() {
           case "enter", "ctrl+c", "q":
               // Exit or transition to another state
               return m, tea.Quit
           }
       }
       return m, cmd
   }
   ```

4. **Update `view.go`**:

   ```go
   // view.go

   func (m model) View() string {
       switch m.state {
       // Existing states
       // ...

       case stateCustomMessageInput:
           return fmt.Sprintf("Enter your custom message:\n\n%s\n\n(Press Enter to continue)", m.customMessageInput.View())

       case stateDisplayCustomMessage:
           return fmt.Sprintf("Your custom message is:\n\n%s\n\n(Press Enter to exit)", m.customMessage)
       }
   }
   ```

---

## General Tips for Adding Features

- **Plan Your State Flow**: Before coding, outline how users will interact with your feature and how states will transition.
- **Consistent Naming**: Use clear and consistent names for states, variables, and functions.
- **Update All Components**: Ensure that any new states or fields are updated across all relevant files.
- **Error Handling**: Add appropriate error handling and user feedback for invalid inputs or unexpected behavior.
- **Testing**: Thoroughly test your new feature to ensure it works as intended within the state machine.

## Understanding BubbleTea's Architecture

BubbleTea applications follow the Model-View-Update (MVU) architecture:

- **Model**: Represents the application's state.
- **View**: A function of the model that returns a string to be rendered.
- **Update**: A function that takes the current model and a message (e.g., user input) and returns an updated model.

By understanding this pattern, you can effectively manage state transitions and user interactions.

## Additional Resources

- **BubbleTea Documentation**: [BubbleTea GitHub Repository](https://github.com/charmbracelet/bubbletea)
- **BubbleTea Examples**: Explore the examples provided in the BubbleTea repository to see practical implementations.
- **Go Concurrency Patterns**: Understanding Go's concurrency model can help when dealing with asynchronous messages.

