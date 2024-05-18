import reactLogo from 'assets/react.svg';
import viteLogo from '/vite.svg';

export default function Home() {
  return (
    <div className="h-dvh w-dvw text-center place-content-center">
      <div className="flex justify-center">
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="h-32 p-6" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="animate-spin h-32 p-6" style={{animation: 'spin 5s linear infinite'}} alt="React logo" />
        </a>
      </div>
      <h1 className="text-5xl font-bold">Vite + React</h1>
      <div className="p-8">
        <p>
          Edit <code>app/pages/home/index.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="text-gray-400">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}