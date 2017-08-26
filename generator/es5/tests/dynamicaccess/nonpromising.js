$module('nonpromising', function () {
  var $static = this;
  this.$struct('5d13a931', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (field) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        field: field,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'field', 'field', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|2549c819<5d13a931>": true,
        "equals|4|2549c819<f361570c>": true,
        "Stringify|2|2549c819<ec87fc3f>": true,
        "Mapping|2|2549c819<95776681<any>>": true,
        "Clone|2|2549c819<5d13a931>": true,
        "String|2|2549c819<ec87fc3f>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.doSomething = function () {
    var ss;
    var ssa;
    ss = $g.nonpromising.SomeStruct.new($t.fastbox(42, $g.________testlib.basictypes.Integer));
    ssa = ss;
    return $t.cast($t.dynamicaccess(ssa, 'field', false), $g.________testlib.basictypes.Integer, false);
  };
  $static.TEST = function () {
    return $t.fastbox($g.nonpromising.doSomething().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
