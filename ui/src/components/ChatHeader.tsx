import React, { useState } from "react";
import { AiOutlineNumber } from 'react-icons/ai'


function Form({ setSocket, socket, setMessages }: { socket: WebSocket | null, setSocket: React.Dispatch<React.SetStateAction<WebSocket | null>>, setMessages: React.Dispatch<React.SetStateAction<string[]>> }) {
    const [message, setMessage] = useState<string>("")
    const [isConnected, setIsConnected] = useState<boolean>(false)

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const formData = new FormData(e.target as HTMLFormElement);
        const data = Object.fromEntries(formData)
        let socket = new WebSocket(`wss://api.chat-app.heyimcarlos.dev/room/${data.room}/user/${data.username}`);
        socket.onopen = () => {
            console.log("connected")
            setSocket(socket)
            setIsConnected(true)
        }
        socket.onclose = () => {
            console.log("disconnected")
            setSocket(null)
            setIsConnected(false)
        }
        socket.onmessage = (e) => {
            console.log("message received: ", e.data)
            setMessages((messages) => [...messages, e.data])
        }

    }


    return (
        <form onSubmit={(e) => {
            console.log("submitting: ", e)
            handleSubmit(e)
        }}>
            <div className="flex justify-between">
                {socket ? (
                    <div className="w-full mr-3">
                        <input
                            onChange={(e) => {
                                e.preventDefault()
                                setMessage(e.target.value)
                            }}
                            value={message}
                            type="text"
                            name="message"
                            id="message"
                            className="block w-full rounded-md text-base px-2 py-2 border-custom-zinc dark:border-zinc-800 dark:bg-custom-black border shadow-sm placeholder:text-zinc-400 sm:leading-6"
                            placeholder="Start chatting!"
                        />
                    </div>
                ) : (
                    <div className="flex gap-2">
                        <div>
                            <input
                                type="name"
                                name="username"
                                id="username"
                                className="block w-full rounded-md text-base px-2 py-2 border-custom-zinc dark:border-zinc-800 dark:bg-custom-black border shadow-sm placeholder:text-zinc-400 sm:leading-6"
                                placeholder="Your username"
                            />
                        </div>

                        <div className="relative w-28">
                            <div className="pointer-events-none absolute inset-y-0 pb-0.5 left-0 flex items-center pl-3">
                                <AiOutlineNumber className="h-5 w-5 text-zinc-500" aria-hidden="true" />
                            </div>
                            <input
                                type="number"
                                name="room"
                                id="room"
                                className="block w-full rounded-md border border-custom-zinc dark:border-zinc-800 dark:bg-custom-black placeholder:text-zinc-400 text-base pr-2 py-2 pl-10 text-custom-zinc dark:text-zinc-400 dark:ring-custom-white shadow-sm sm:leading-6"
                                placeholder="room"
                            />
                        </div>
                    </div>
                )}
                {isConnected ? (
                    <div className="flex gap-2">
                        <button
                            onClick={(e) => {
                                e.preventDefault()
                                socket && socket.send(message)
                                setMessage("")
                            }}
                            className="rounded bg-green-500 px-2 py-1 text-base font-semibold text-white shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >
                            Send
                        </button>
                        <button
                            onClick={(e) => {
                                e.preventDefault()
                                socket && socket.close()
                                setSocket(null)
                            }}
                            className="rounded bg-red-500 px-2 py-1 text-base font-semibold text-white shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >
                            Disconnect
                        </button>

                    </div>

                ) : (
                    <button
                        type="submit"
                        className="rounded bg-custom-black dark:bg-zinc-700 px-2 py-1 text-base font-semibold text-white shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                    >
                        Connect
                    </button>

                )}
            </div>
        </form >
    );
}

export default function ChatHeader() {
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const [messages, setMessages] = useState<string[]>([])

    return (
        <div className="m-4 space-y-3 dark:text-zinc-400">
            <div className="gap-2">
                <div className="flex w-full justify-between items-center gap-2">
                    {socket ? (
                        <div className="flex justify-end w-full items-center">
                            <div className={'flex-none rounded-full p-1 text-green-400 bg-green-400/10'}>
                                <div className="h-1.5 w-1.5 rounded-full bg-current" />
                            </div>
                            Online
                        </div>
                    ) : (
                        <>
                            <h3 className="text-lg pl-0.5 font-medium leading-6 text-custom-zinc dark:text-zinc-300">Enter a username, room number and click connect</h3>
                            <div className="flex items-center gap-2">
                                <div className={'flex-none rounded-full p-1 text-red-400 bg-red-400/10'}>
                                    <div className="h-1.5 w-1.5 rounded-full bg-current" />
                                </div>
                                Offline
                            </div>
                        </>
                    )}
                </div>
            </div>
            <Form setMessages={setMessages} setSocket={setSocket} socket={socket} />
            <div className="mt-5 h-[450px] relative">

                <pre className="relative h-full text-custom-zinc dark:text-zinc-500 bg-zinc-50 dark:bg-zinc-200 p-4 overflow-y-scroll">
                    {socket ? (
                        <>
                            {
                                messages.map((message, idx) => {
                                    if (message.startsWith("you just joined room: ")) {
                                        return (
                                            <p key={idx} className="text-green-500">{message}</p>
                                        )
                                    }
                                    return (
                                        <p key={idx}>{message}</p>
                                    )

                                })
                            }
                        </>
                    ) : (
                        <div>Please connect to a room to start chatting...</div>
                    )}
                </pre>
            </div>
        </div >
    )
}

