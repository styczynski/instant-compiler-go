.class public <input>
.super java/lang/Object
   .method <init>()V
      .limit stack 1
      .limit locals 1
      aload_0
      invokespecial java/lang/Object/<init>()V
      return
   .end method
   .method public static main([Ljava/lang/String;)V
      .limit stack 2
      .limit locals 1
      getstatic java/lang/System/out Ljava/io/PrintStream;
      iconst_2
      invokevirtual java/io/PrintStream/println(I)V
      return
   .end method