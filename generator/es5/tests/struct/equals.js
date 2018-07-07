$module('equals', function () {
  var $static = this;
  this.$struct('9285bf36', 'Foo', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeValue, AnotherValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeValue: SomeValue,
        AnotherValue: AnotherValue,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeValue', 'SomeValue', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherValue', 'AnotherValue', function () {
      return $g.equals.Bar;
    }, function () {
      return $g.equals.Bar;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|6caba86c<9285bf36>": true,
        "equals|4|6caba86c<0e92a8bc>": true,
        "Stringify|2|6caba86c<e38ac9b0>": true,
        "Mapping|2|6caba86c<c518fe3b<any>>": true,
        "Clone|2|6caba86c<9285bf36>": true,
        "String|2|6caba86c<e38ac9b0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('4fb96a52', 'Bar', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (StringValue) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        StringValue: StringValue,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'StringValue', 'StringValue', function () {
      return $g.________testlib.basictypes.String;
    }, function () {
      return $g.________testlib.basictypes.String;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|6caba86c<4fb96a52>": true,
        "equals|4|6caba86c<0e92a8bc>": true,
        "Stringify|2|6caba86c<e38ac9b0>": true,
        "Mapping|2|6caba86c<c518fe3b<any>>": true,
        "Clone|2|6caba86c<4fb96a52>": true,
        "String|2|6caba86c<e38ac9b0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var copy;
    var different;
    var first;
    var second;
    first = $g.equals.Foo.new($t.fastbox(42, $g.________testlib.basictypes.Integer), $g.equals.Bar.new($t.fastbox('hello world', $g.________testlib.basictypes.String)));
    second = first;
    copy = $g.equals.Foo.new($t.fastbox(42, $g.________testlib.basictypes.Integer), $g.equals.Bar.new($t.fastbox('hello world', $g.________testlib.basictypes.String)));
    different = $g.equals.Foo.new($t.fastbox(42, $g.________testlib.basictypes.Integer), $g.equals.Bar.new($t.fastbox('hello worlds!', $g.________testlib.basictypes.String)));
    return $t.fastbox((($g.equals.Foo.$equals(first, second).$wrapped && $g.equals.Foo.$equals(first, copy).$wrapped) && !$g.equals.Foo.$equals(first, different).$wrapped) && !$g.equals.Foo.$equals(copy, different).$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
