import { useEffect, useState } from "react"
import style from './table.module.css'

const Table: React.FC = () => {
    const [data, setData] = useState<{ [key: string]: string }>({});
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:8000/api/v1/getAllData'); // Replace with your API endpoint
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const jsonData = await response.json();
                if (jsonData.error === 1) {
                    throw new Error("Couldn't Retrieve data")
                }
                setData(jsonData.data);
            } catch (e: any) {
                setError(e.message);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className={style.tableContainer}>
            <table className={style.styledTable}>
                <thead>
                    <tr>
                        <th>Shorten URL</th>
                        <th>Main URL</th>
                    </tr>
                </thead>
                <tbody>
                    {Object.entries(data).map(([key, val]) => (
                        <tr key={key}>
                            <td><center>{key}</center></td>
                            <td><center>{val}</center></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default Table