export interface SpinResult {
    participantIds: number[],
    winnerId: number
}

export interface SpinResultResponse {
    result: SpinResult,
    error: string
}