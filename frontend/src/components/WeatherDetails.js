import axios from "axios"
import Cookies from "js-cookie"
import React, { useEffect, useState } from "react"
import { Link, useParams } from "react-router-dom"
import "./WeatherDetails.css" // Импорт файла стилей

const WeatherDetails = () => {
	const { id } = useParams()
	const [observation, setObservation] = useState(null)

	var userRole = Cookies.get("X-User-Role")

	useEffect(() => {
		axios
			.get(`http://localhost:8080/weather/${id}`)
			.then(response => setObservation(response.data))
			.catch(error => console.error(error))
	}, [id])

	const formatDate = timestamp => {
		const date = new Date(timestamp)
		const day = String(date.getDate()).padStart(2, "0") // Добавляем ведущий ноль
		const month = String(date.getMonth() + 1).padStart(2, "0") // Месяц начинается с 0
		const year = date.getFullYear()
		return `${day}.${month}.${year}`
	}

	if (!observation) return <div className="loading">Loading...</div>

	return (
		<div className="details-container">
			<h1 className="details-title">Weather Details</h1>
			<div className="details-info">
				<p>
					<strong>Date:</strong> {formatDate(observation.timestamp)}
				</p>
				<p>
					<strong>Temperature:</strong> {observation.temperature}°C
				</p>
				<p>
					<strong>Humidity:</strong> {observation.humidity}%
				</p>
				<p>
					<strong>Pressure:</strong> {observation.pressure} hPa
				</p>
				<p>
					<strong>Wind Speed:</strong> {observation.wind_speed} m/s
				</p>
				<p>
					<strong>Weather Status:</strong> {observation.weather_status}
				</p>
			</div>
			<div className="details-actions">
				<Link to="/" className="back-button">
					Back to List
				</Link>
				{userRole === "admin" && (
					<Link to={`/update/${id}`} className="edit-button">
						Edit Observation
					</Link>
				)}
			</div>
		</div>
	)
}

export default WeatherDetails
