export interface Cities {
    id: string;
    name: string;
    latitude: number;
    longitude: number;
}

export interface WeatherData {
    city_id: string;
    city_name: string;
    temperature: number;
    humidity: number;
    wind_speed: number;
    condition: string;
}

export interface CompareResult {
    cities: WeatherData[];
    summary: WeatherSummary;
    error: null | string;
}

export interface WeatherSummary {
    average_temperature: number;
    average_humidity: number;
    average_wind_speed: number;
    hotter_city: string;
    colder_city: string;
    windy_city: string;
    ranking: string[];
    by_condition: Record<string, string[]>;
}