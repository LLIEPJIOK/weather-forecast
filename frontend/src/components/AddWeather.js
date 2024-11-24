import axios from "axios"
import React, { useState } from "react"
import AddWeatherForm from "./AddWeatherForm"
import observationTemplate from "./observationTemplate" // Импорт шаблона

const AddWeather = () => {
	const [observation, setObservation] = useState(observationTemplate) // Используем шаблон

	const handleChange = e => {
		const { name, value } = e.target

		// Directly update properties if they exist in the observation template
		if (name in observation) {
			setObservation({
				...observation,
				[name]: value, // Update the corresponding property directly
			})
		}
	}

	const handleSubmit = e => {
		e.preventDefault()
		const formattedObservation = {
			...observation,
			timestamp: new Date(observation.timestamp).toISOString(), // Приводим к RFC3339
			temperature: parseFloat(observation.temperature), // Преобразуем в число
			humidity: parseFloat(observation.humidity), // Преобразуем в число
			pressure: parseFloat(observation.pressure), // Преобразуем в число
			wind_speed: parseFloat(observation.windSpeed), // Преобразуем в число
			weather_status: observation.weatherStatus,
		}
		axios
			.post("http://localhost:8080/weather", formattedObservation)
			.then(() => {
				setObservation(observationTemplate)
				console.log("Observation added successfully!")
			})
			.catch(error => console.error("Error adding observation:", error))
	}

	return (
		<AddWeatherForm
			observation={observation}
			handleChange={handleChange}
			handleSubmit={handleSubmit}
		/>
	)
}

export default AddWeather
