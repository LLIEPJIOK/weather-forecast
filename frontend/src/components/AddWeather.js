import axios from "axios"
import Cookies from "js-cookie"
import React, { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import AddWeatherForm from "./AddWeatherForm"
import observationTemplate from "./observationTemplate"

const AddWeather = () => {
	const [observation, setObservation] = useState(observationTemplate)
	const navigate = useNavigate()

	useEffect(() => {
		var userRole = Cookies.get("X-User-Role")
		if (userRole !== "admin") {
			navigate("/login")
		}
	}, [])

	const handleChange = e => {
		const { name, value } = e.target

		if (name in observation) {
			setObservation({
				...observation,
				[name]: value,
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
