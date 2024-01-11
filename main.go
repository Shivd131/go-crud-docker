package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// STRUCTURES:
// Task structure represents a task within a project
type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// Project structure represents a project with a list of tasks
type Project struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Tasks []Task `json:"tasks"`
}

var projects []Project

func main() {
	router := gin.Default()

	// Routes
	router.GET("/projects", getProjects)
	router.GET("/projects/:id", getProject)
	router.POST("/projects", createProject)
	router.PUT("/projects/:id", updateProject)
	router.DELETE("/projects/:id", deleteProject)

	router.POST("/projects/:id/tasks", createTask)
	router.PUT("/projects/:id/tasks/:taskID", updateTask)
	router.DELETE("/projects/:id/tasks/:taskID", deleteTask)

	router.Run(":8080")
}

// returns the list of all projects
func getProjects(c *gin.Context) {
	c.JSON(http.StatusOK, projects)
}

// returns a specific project by ID
func getProject(c *gin.Context) {
	id := getIDParam(c)
	project := findProjectByID(id)
	if project.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

// creates a new project
func createProject(c *gin.Context) {
	var newProject Project
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProject.ID = generateID()
	projects = append(projects, newProject)
	c.JSON(http.StatusCreated, newProject)
}

// updates an existing project by ID
func updateProject(c *gin.Context) {
	id := getIDParam(c)
	index := findProjectIndexByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var updatedProject Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProject.ID = id
	projects[index] = updatedProject
	c.JSON(http.StatusOK, updatedProject)
}

// deletes a project by ID
func deleteProject(c *gin.Context) {
	id := getIDParam(c)
	index := findProjectIndexByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	projects = append(projects[:index], projects[index+1:]...)
	c.JSON(http.StatusNoContent, nil)
}

// adds a new task to a project
func createTask(c *gin.Context) {
	projectID := getIDParam(c)
	index := findProjectIndexByID(projectID)
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.ID = generateID()
	projects[index].Tasks = append(projects[index].Tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

// updates an existing task within a project by ID
func updateTask(c *gin.Context) {
	projectID := getIDParam(c)
	taskID := getTaskIDParam(c)

	projectIndex := findProjectIndexByID(projectID)
	if projectIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	taskIndex := findTaskIndexByID(projectIndex, taskID)
	if taskIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask.ID = taskID
	projects[projectIndex].Tasks[taskIndex] = updatedTask
	c.JSON(http.StatusOK, updatedTask)
}

// deletes a task within a project by ID
func deleteTask(c *gin.Context) {
	projectID := getIDParam(c)
	taskID := getTaskIDParam(c)

	projectIndex := findProjectIndexByID(projectID)
	if projectIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	taskIndex := findTaskIndexByID(projectIndex, taskID)
	if taskIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	projects[projectIndex].Tasks = append(projects[projectIndex].Tasks[:taskIndex], projects[projectIndex].Tasks[taskIndex+1:]...)
	c.JSON(http.StatusNoContent, nil)
}

// get project ID from the URL parameters
func getIDParam(c *gin.Context) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0
	}
	return id
}

// get the task ID from the URL parameters
func getTaskIDParam(c *gin.Context) int {
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		return 0
	}
	return taskID
}

// to find a project by ID
func findProjectByID(id int) Project {
	for _, project := range projects {
		if project.ID == id {
			return project
		}
	}
	return Project{}
}

// to find the index of a project by ID
func findProjectIndexByID(id int) int {
	for i, project := range projects {
		if project.ID == id {
			return i
		}
	}
	return -1
}

// to find the index of a task within a project by ID
func findTaskIndexByID(projectIndex, taskID int) int {
	for i, task := range projects[projectIndex].Tasks {
		if task.ID == taskID {
			return i
		}
	}
	return -1
}

// to generate a new unique ID for a project or task
func generateID() int {

	return len(projects) + 1
}
