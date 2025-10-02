'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Search, Plus, MoreHorizontal, Phone, Video, Users } from 'lucide-react'

interface Message {
  id: number
  text: string
  sender: 'user' | 'system'
  timestamp: Date
}

export default function Home() {
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [isMounted, setIsMounted] = useState(false)

  useEffect(() => {
    setIsMounted(true)
    // 模拟初始消息
    setMessages([
      { id: 1, text: '你好！欢迎使用聊天应用', sender: 'system', timestamp: new Date() },
      { id: 2, text: '这是一个类似钉钉的聊天界面', sender: 'system', timestamp: new Date(Date.now() - 300000) },
    ])
  }, [])

  const handleSendMessage = () => {
    if (newMessage.trim() && isMounted) {
      setMessages([
        ...messages,
        {
          id: messages.length + 1,
          text: newMessage,
          sender: 'user',
          timestamp: new Date(),
        },
      ])
      setNewMessage('')
    }
  }

  const formatTime = (date: Date) => {
    if (!isMounted) return ''
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  if (!isMounted) {
    return (
      <div className="flex h-screen bg-background items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-muted-foreground">加载中...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-background">
      {/* 侧边栏 */}
      <div className="w-80 border-r bg-background">
        {/* 顶部搜索和新建 */}
        <div className="p-4 border-b">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-xl font-bold text-primary">ChatApp</h1>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <Plus className="h-4 w-4" />
            </Button>
          </div>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="搜索聊天室或联系人..."
              className="pl-10 bg-muted/50 border-none"
            />
          </div>
        </div>
        
        {/* 聊天室列表 */}
        <ScrollArea className="h-[calc(100vh-140px)]">
          <div className="p-2">
            <div className="space-y-1">
              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors bg-accent/50">
                <Avatar className="h-10 w-10">
                  <AvatarFallback className="bg-primary text-primary-foreground">
                    <Users className="h-5 w-5" />
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-semibold">技术讨论</p>
                    <span className="text-xs text-muted-foreground">{formatTime(new Date())}</span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">张三: 大家看看这个新功能...</p>
                </div>
                <Badge variant="secondary" className="h-5">3</Badge>
              </div>

              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors">
                <Avatar className="h-10 w-10">
                  <AvatarFallback className="bg-blue-500 text-white">D</AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-semibold">设计团队</p>
                    <span className="text-xs text-muted-foreground">昨天</span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">设计稿已更新</p>
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors">
                <Avatar className="h-10 w-10">
                  <AvatarFallback className="bg-green-500 text-white">P</AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-semibold">产品规划</p>
                    <span className="text-xs text-muted-foreground">周三</span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">新版本需求讨论</p>
                </div>
              </div>
            </div>

            <Separator className="my-4" />

            {/* 联系人列表 */}
            <div className="space-y-1">
              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors">
                <div className="relative">
                  <Avatar className="h-10 w-10">
                    <AvatarFallback className="bg-blue-500 text-white">张</AvatarFallback>
                  </Avatar>
                  <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-background"></div>
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium">张三</p>
                  <p className="text-xs text-muted-foreground">产品经理</p>
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors">
                <div className="relative">
                  <Avatar className="h-10 w-10">
                    <AvatarFallback className="bg-green-500 text-white">李</AvatarFallback>
                  </Avatar>
                  <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-yellow-500 rounded-full border-2 border-background"></div>
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium">李四</p>
                  <p className="text-xs text-muted-foreground">UI设计师</p>
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 rounded-lg hover:bg-accent cursor-pointer transition-colors">
                <div className="relative">
                  <Avatar className="h-10 w-10">
                    <AvatarFallback className="bg-purple-500 text-white">王</AvatarFallback>
                  </Avatar>
                  <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-gray-400 rounded-full border-2 border-background"></div>
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium">王五</p>
                  <p className="text-xs text-muted-foreground">后端开发</p>
                </div>
              </div>
            </div>
          </div>
        </ScrollArea>
      </div>

      {/* 主聊天区域 */}
      <div className="flex-1 flex flex-col">
        {/* 聊天头部 */}
        <div className="p-4 border-b">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarFallback className="bg-primary text-primary-foreground">
                  <Users className="h-5 w-5" />
                </AvatarFallback>
              </Avatar>
              <div>
                <h2 className="font-semibold text-lg">技术讨论</h2>
                <p className="text-sm text-muted-foreground">张三、李四、王五 在线</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="icon" className="h-9 w-9">
                <Phone className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-9 w-9">
                <Video className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-9 w-9">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>

        {/* 消息区域 */}
        <ScrollArea className="flex-1">
          <div className="p-6 space-y-6">
            {/* 欢迎消息 */}
            <div className="text-center">
              <div className="inline-flex items-center gap-2 px-4 py-2 bg-muted rounded-full">
                <Avatar className="h-6 w-6">
                  <AvatarFallback className="bg-primary text-primary-foreground text-xs">系</AvatarFallback>
                </Avatar>
                <span className="text-sm text-muted-foreground">今天 09:30</span>
              </div>
              <p className="text-sm text-muted-foreground mt-2">这是聊天开始的地方，发送消息开始对话</p>
            </div>

            {/* 消息列表 */}
            {messages.map((message) => (
              <div
                key={message.id}
                className={`flex gap-3 ${
                  message.sender === 'user' ? 'flex-row-reverse' : ''
                }`}
              >
                <Avatar className="h-8 w-8 flex-shrink-0">
                  <AvatarFallback className={
                    message.sender === 'user' 
                      ? 'bg-primary text-primary-foreground' 
                      : 'bg-muted text-muted-foreground'
                  }>
                    {message.sender === 'user' ? '我' : '系'}
                  </AvatarFallback>
                </Avatar>
                <div className={`max-w-[70%] ${
                  message.sender === 'user' ? 'text-right' : ''
                }`}>
                  <div className={
                    message.sender === 'user' 
                      ? 'bg-primary text-primary-foreground rounded-2xl rounded-tr-md px-4 py-2'
                      : 'bg-muted rounded-2xl rounded-tl-md px-4 py-2'
                  }>
                    <p className="text-sm leading-relaxed">{message.text}</p>
                  </div>
                  <p className="text-xs text-muted-foreground mt-1 px-1">
                    {formatTime(message.timestamp)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>

        {/* 输入区域 */}
        <div className="p-4 border-t">
          <div className="flex gap-2">
            <Button variant="ghost" size="icon" className="h-10 w-10">
              <Plus className="h-5 w-5" />
            </Button>
            <Input
              placeholder="输入消息..."
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
              className="flex-1 rounded-full px-4"
            />
            <Button 
              onClick={handleSendMessage} 
              className="rounded-full bg-primary hover:bg-primary/90 px-6"
              disabled={!newMessage.trim()}
            >
              发送
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
