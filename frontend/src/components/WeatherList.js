import axios from "axios"
import React, { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import "./WeatherList.css" // Импорт стилей

const WeatherList = () => {
	const [weatherObservations, setWeatherObservations] = useState([])

	useEffect(() => {
		axios
			.get("http://localhost:8080/weathers")
			.then(response => setWeatherObservations(response.data))
			.catch(error => console.error("Error fetching data:", error))
	}, [])

	// Функция для форматирования даты в dd.mm.yyyy
	const formatDate = timestamp => {
		const date = new Date(timestamp)
		const day = String(date.getDate()).padStart(2, "0") // Добавляем ведущий ноль
		const month = String(date.getMonth() + 1).padStart(2, "0") // Месяц начинается с 0
		const year = date.getFullYear()
		return `${day}.${month}.${year}`
	}

	const truncateText = (text, length) => {
		if (text.length > length) {
			return text.substring(0, length) + "..."
		}
		return text
	}

	return (
		<div className="weather-container">
			<h1 className="weather-title">Weather Observations</h1>
			<ul className="weather-list">
				{weatherObservations.map(ob => (
					<li key={ob.id} className="weather-item">
						<div className="weather-info">
							<span className="weather-date">{formatDate(ob.timestamp)}</span>
							<span className="weather-status">
								{truncateText(ob.weather_status, 25)}
							</span>
						</div>
						<div className="details-button-container">
							<Link to={`/details/${ob.id}`} className="details-button">
								Details
							</Link>
						</div>
					</li>
				))}
			</ul>
		</div>
	)
}

export default WeatherList
