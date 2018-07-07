$module('basic', function () {
  var $static = this;
  this.$struct('a76166f4', 'AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherBool', 'AnotherBool', function () {
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|6caba86c<a76166f4>": true,
        "equals|4|6caba86c<0e92a8bc>": true,
        "Stringify|2|6caba86c<e38ac9b0>": true,
        "Mapping|2|6caba86c<c518fe3b<any>>": true,
        "Clone|2|6caba86c<a76166f4>": true,
        "String|2|6caba86c<e38ac9b0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('1a1b7840', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
        SomeInstance: SomeInstance,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    $t.defineStructField($static, 'SomeInstance', 'SomeInstance', function () {
      return $g.basic.AnotherStruct;
    }, function () {
      return $g.basic.AnotherStruct;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|6caba86c<1a1b7840>": true,
        "equals|4|6caba86c<0e92a8bc>": true,
        "Stringify|2|6caba86c<e38ac9b0>": true,
        "Mapping|2|6caba86c<c518fe3b<any>>": true,
        "Clone|2|6caba86c<1a1b7840>": true,
        "String|2|6caba86c<e38ac9b0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var ss;
    ss = $g.basic.SomeStruct.new($t.fastbox(42, $g.________testlib.basictypes.Integer), $t.fastbox(true, $g.________testlib.basictypes.Boolean), $g.basic.AnotherStruct.new($t.fastbox(true, $g.________testlib.basictypes.Boolean)));
    return $t.fastbox(((ss.SomeField.$wrapped == 42) && ss.AnotherField.$wrapped) && ss.SomeInstance.AnotherBool.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
