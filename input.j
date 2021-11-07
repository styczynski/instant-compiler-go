.class public input
.super java/lang/Object
   .method <init>()V
      .limit stack 1
      .limit locals 1
      aload_0
      invokespecial java/lang/Object/<init>()V
      return
   .end method
   .method public static main([Ljava/lang/String;)V
      .limit stack 3
      .limit locals 1
      getstatic java/lang/System/out Ljava/io/PrintStream;
      iconst_4
      iconst_5
      iadd
      invokevirtual java/io/PrintStream/println(I)V
      return
   .end method