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
      .limit locals 2
      getstatic java/lang/System/out Ljava/io/PrintStream;
      iconst_2
      iconst_2
      iadd
      invokevirtual java/io/PrintStream/println(I)V
      iconst_2
      istore_1
      iconst_1
      iconst_4
      swap
      idiv
      iconst_1
      swap
      idiv
      istore_1
      return
   .end method